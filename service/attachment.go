package service

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/utils"
	"gin-web/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type AttachmentService struct {
	*BasicService
	repository *repository.AttachmentRepository
}

func NewAttachmentService(ctx *gin.Context) *AttachmentService {
	return &AttachmentService{
		BasicService: NewBasicService(ctx),
		repository:   repository.NewAttachmentRepository(ctx),
	}
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

// 批量保存多个文件
func (a *AttachmentService) SaveFiles(files []*multipart.FileHeader, uid uint) ([]*models.Attachment, error) {
	var info = make([]*AttachmentInfo, 0, len(files))
	projectDir, temp := os.Getwd()
	if temp != nil {
		return nil, temp
	}
	targetDir := viper.GetString("system.uploadAttachmentDir")
	var uploadAttachmentDir = filepath.Join(targetDir, time.Now().Format(time.DateOnly))
	if strings.TrimSpace(targetDir) == "" {
		uploadAttachmentDir = filepath.Join(projectDir, "upload", time.Now().Format(time.DateOnly))
	}
	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer func(f multipart.File) {
			_ = f.Close()
		}(f)
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
		defer func(t *os.File) {
			_ = t.Close()
		}(tempFile)
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
		info = append(info, &AttachmentInfo{
			uniqueName: fmt.Sprintf("%s%s", fileHash, fileMd5),
			classify:   getFileType(buffer.Bytes()),
		})
		targetFilePath := filepath.Join(uploadAttachmentDir, fmt.Sprintf("%s%s%s", fileHash, fileMd5, ext))
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
			Name: item.Filename,
			Uid:  uid,
			Type: info[index].classify,
			Size: item.Size,
			Hash: info[index].uniqueName,
		}
	})
	return a.repository.Create(a.ctx.Request.Context(), attachments)
}
