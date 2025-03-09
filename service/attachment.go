package service

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/response"
	"gin-web/pkg/utils"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/samber/lo"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type AttachmentRepository interface {
	Create(ctx context.Context, records []*models.Attachment) error
}

type AttachmentService struct {
	*BasicService
	repository AttachmentRepository
}

var (
	attachmentOnce    sync.Once
	attachmentService *AttachmentService
)

func NewAttachmentService(basic *BasicService, repo AttachmentRepository) *AttachmentService {
	attachmentOnce.Do(func() {
		attachmentService = &AttachmentService{
			BasicService: basic,
			repository:   repo,
		}
	})
	return attachmentService
}

// https://github.com/h2non/filetype
const (
	other = iota + 1
	image
	video
	audio
	archive
	documents
	font
	application
)

type AttachmentInfo struct {
	uniqueName string
	classify   int8
	path       string
}

func getFileType(header []byte) int8 {
	switch {
	case filetype.IsImage(header):
		return image
	case filetype.IsVideo(header):
		return video
	case filetype.IsAudio(header):
		return audio
	case filetype.IsArchive(header):
		return archive
	case filetype.IsFont(header):
		return font
	case filetype.IsApplication(header):
		return application
	case filetype.IsDocument(header):
		return documents
	default:
		return other
	}
}

func (a *AttachmentService) SaveFiles(ctx context.Context, files []*multipart.FileHeader, uid uint) ([]*response.UploadAttachmentResponse, error) {
	var info = make([]*AttachmentInfo, 0, len(files))
	projectDir, temp := os.Getwd()
	if temp != nil {
		return nil, temp
	}
	targetDir := a.viper.GetString("system.uploadAttachmentDir")
	var uploadAttachmentDir = filepath.Join(targetDir, time.Now().Format(time.DateOnly))
	if strings.TrimSpace(targetDir) == "" {
		uploadAttachmentDir = filepath.Join(projectDir, "upload", time.Now().Format(time.DateOnly))
	}
	var openFiles = make([]multipart.File, 0)
	var createFiles = make([]*os.File, 0)
	defer func() {
		for _, f := range openFiles {
			_ = f.Close()
		}
		for _, f := range createFiles {
			_ = f.Close()
		}
	}()
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			return nil, err
		}
		openFiles = append(openFiles, f)
		ext := filepath.Ext(file.Filename)
		// 同时将文件内容写入hash和上传文件夹中的临时文件
		tempFileName := fmt.Sprintf("temp_%s", uuid.New().String())
		tempFilePath := filepath.Join(uploadAttachmentDir, fmt.Sprintf("%s%s", tempFileName, ext))
		if err = os.MkdirAll(filepath.Dir(tempFilePath), 0750); err != nil {
			return nil, err
		}
		tempFile, err := os.Create(tempFilePath)
		if err != nil {
			return nil, err
		}
		createFiles = append(createFiles, tempFile)
		// 获取前262个字节用于判断文件类型
		buffer := bytes.NewBuffer(make([]byte, 0, 262))
		_, err = io.CopyN(buffer, f, 262)
		if err != nil {
			return nil, err
		}
		// 还要继续处理文件，重置文件指针
		if _, seekErr := f.Seek(0, io.SeekStart); seekErr != nil {
			return nil, seekErr
		}
		hash := sha1.New()
		md := md5.New()
		if _, copyErr := io.Copy(io.MultiWriter(hash, md, tempFile), f); copyErr != nil {
			return nil, err
		}
		fileHash := hex.EncodeToString(hash.Sum(nil))
		fileMd5 := hex.EncodeToString(md.Sum(nil))
		targetFilePath := filepath.Join(uploadAttachmentDir, fmt.Sprintf("%s%s%s", fileHash, fileMd5, ext))
		info = append(info, &AttachmentInfo{
			uniqueName: fmt.Sprintf("%s%s", fileHash, fileMd5),
			classify:   getFileType(buffer.Bytes()),
			path:       targetFilePath,
		})
		// 文件如果不存在才创建,避免重复创建多个内容一样的文件
		exists, err := utils.PathExists(targetFilePath)
		if err != nil {
			return nil, err
		}
		// 文件不存在才创建
		if !(exists && utils.IsFile(targetFilePath)) {
			err = os.Rename(tempFilePath, targetFilePath)
			if err != nil {
				return nil, err
			}
		} else {
			delErr := os.Remove(tempFilePath)
			if delErr != nil {
				return nil, delErr
			}
		}
	}
	attachments := lo.Map(files, func(item *multipart.FileHeader, index int) *models.Attachment {
		return &models.Attachment{
			Name:      item.Filename,
			CreatorId: &uid,
			Type:      info[index].classify,
			Size:      item.Size,
			Hash:      info[index].uniqueName,
			Path:      info[index].path,
		}
	})
	// TODO:创建一个事务，创建文件记录的同时为文件授权
	temp = a.repository.Create(ctx, attachments)
	if temp != nil {
		return nil, temp
	}
	ret := lo.Map(attachments, func(item *models.Attachment, _ int) *response.UploadAttachmentResponse {
		return &response.UploadAttachmentResponse{
			ID:   item.ID,
			Name: item.Name,
			Size: item.Size,
			Path: filepath.Join(string(filepath.Separator), item.Path), // TODO:返回预览/下载接口完整路径
		}
	})
	return ret, nil
}

func (a *AttachmentService) SaveFile(ctx context.Context, file *multipart.FileHeader, uid uint) (*response.UploadAttachmentResponse, error) {
	files, err := a.SaveFiles(ctx, []*multipart.FileHeader{file}, uid)
	if err != nil || files == nil || len(files) == 0 {
		return nil, err
	}
	return files[0], err
}
