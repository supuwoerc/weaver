package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestPathExists(t *testing.T) {
	dir := t.TempDir() // 会自动清理的零食目录
	tempFile := filepath.Join(dir, "temp")
	if err := os.WriteFile(tempFile, []byte("test file content"), 0644); err != nil {
		t.Error(err)
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "EmptyPath",
			args: args{
				path: "",
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ExistingDirectory",
			args: args{
				path: dir,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "ExistingFile",
			args: args{
				path: tempFile,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "NonExistentPath",
			args: args{
				path: filepath.Join(dir, "non-existent"),
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PathExists(tt.args.path)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equalf(t, tt.want, got, "PathExists,%v,params:%v", tt.name, tt.args.path)
		})
	}
	// 软链接
	t.Run("Symlink", func(t *testing.T) {
		link := filepath.Join(dir, "symlink")
		err := os.Symlink(tempFile, link) // 创建软链接
		if err != nil {
			t.Error(err)
			return
		}
		exists, err := PathExists(link)
		assert.True(t, exists)
		assert.NoError(t, err)
	})
	t.Run("PermissionDenied", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("Test skipped for root user")
		}
		// 创建父目录并设置无权限
		parentDir := filepath.Join(dir, "no-access")
		require.NoError(t, os.Mkdir(parentDir, 0700)) // 确保父目录可访问以便后续操作
		defer func() {
			require.NoError(t, os.RemoveAll(parentDir)) // 清理父目录
		}()
		// 在父目录下创建子目录，并移除父目录的权限
		targetDir := filepath.Join(parentDir, "target")
		require.NoError(t, os.Mkdir(targetDir, 0700))
		require.NoError(t, os.Chmod(parentDir, 0000)) // 阻断父目录访问
		defer func() {
			require.NoError(t, os.Chmod(parentDir, 0700)) // 恢复权限以便清理
		}()
		exists, err := PathExists(targetDir)
		assert.False(t, exists)
		assert.Error(t, err)
		assert.True(t, os.IsPermission(err), "err not as permission denied")
	})
}

func TestIsDir(t *testing.T) {
	tmpDir := t.TempDir()
	dirPath := filepath.Join(tmpDir, "test_dir")
	filePath := filepath.Join(tmpDir, "test_file")
	require.NoError(t, os.Mkdir(dirPath, 0755))
	require.NoError(t, os.WriteFile(filePath, []byte{}, 0644))
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "ExistingDirectory",
			args: args{
				path: dirPath,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "ExistingFile",
			args: args{
				path: filePath,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "NonExistingPath",
			args: args{
				path: filepath.Join("non-existing-path"),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsDir(tt.args.path)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, os.IsNotExist(err), "should be not exist error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestIsFile(t *testing.T) {
	tmpDir := t.TempDir()
	dirPath := filepath.Join(tmpDir, "test_dir")
	filePath := filepath.Join(tmpDir, "test_file")
	require.NoError(t, os.Mkdir(dirPath, 0755))
	require.NoError(t, os.WriteFile(filePath, []byte{}, 0644))
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "ExistingDirectory",
			args: args{
				path: dirPath,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "ExistingFile",
			args: args{
				path: filePath,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "NonExistingPath",
			args: args{
				path: filepath.Join("non-existing-path"),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsFile(tt.args.path)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, os.IsNotExist(err), "should be not exist error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
