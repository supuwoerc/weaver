package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supuwoerc/weaver/pkg/constant"
)

func TestOSSConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *OSSConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid AWS S3 config",
			config: &OSSConfig{
				Type:            constant.AWSS3,
				Region:          "us-east-1",
				AccessKeyID:     "test-access-key",
				SecretAccessKey: "test-secret-key",
			},
			wantErr: false,
		},
		{
			name: "valid MinIO config",
			config: &OSSConfig{
				Type:            constant.MinIO,
				Endpoint:        "http://localhost:9000",
				Region:          "us-east-1",
				AccessKeyID:     "test-access-key",
				SecretAccessKey: "test-secret-key",
			},
			wantErr: false,
		},
		{
			name: "valid Aliyun OSS config",
			config: &OSSConfig{
				Type:            constant.AliyunOSS,
				Endpoint:        "https://oss-cn-hangzhou.aliyuncs.com",
				Region:          "cn-hangzhou",
				AccessKeyID:     "test-access-key",
				SecretAccessKey: "test-secret-key",
			},
			wantErr: false,
		},
		{
			name: "valid Tencent COS config",
			config: &OSSConfig{
				Type:            constant.TencentCOS,
				Endpoint:        "https://cos.ap-beijing.myqcloud.com",
				Region:          "ap-beijing",
				AccessKeyID:     "test-access-key",
				SecretAccessKey: "test-secret-key",
			},
			wantErr: false,
		},
		{
			name: "missing type",
			config: &OSSConfig{
				Region:          "us-east-1",
				AccessKeyID:     "test-access-key",
				SecretAccessKey: "test-secret-key",
			},
			wantErr: true,
			errMsg:  "storage type is required",
		},
		{
			name: "missing access key",
			config: &OSSConfig{
				Type:            constant.AWSS3,
				Region:          "us-east-1",
				SecretAccessKey: "test-secret-key",
			},
			wantErr: true,
			errMsg:  "access key ID is required",
		},
		{
			name: "missing secret key",
			config: &OSSConfig{
				Type:        constant.AWSS3,
				Region:      "us-east-1",
				AccessKeyID: "test-access-key",
			},
			wantErr: true,
			errMsg:  "secret access key is required",
		},
		{
			name: "missing endpoint for MinIO",
			config: &OSSConfig{
				Type:            constant.MinIO,
				Region:          "us-east-1",
				AccessKeyID:     "test-access-key",
				SecretAccessKey: "test-secret-key",
			},
			wantErr: true,
			errMsg:  "endpoint is required for non-AWS S3 services",
		},
		{
			name: "missing endpoint for Aliyun OSS",
			config: &OSSConfig{
				Type:            constant.AliyunOSS,
				Region:          "cn-hangzhou",
				AccessKeyID:     "test-access-key",
				SecretAccessKey: "test-secret-key",
			},
			wantErr: true,
			errMsg:  "endpoint is required for non-AWS S3 services",
		},
		{
			name: "missing endpoint for Tencent COS",
			config: &OSSConfig{
				Type:            constant.TencentCOS,
				Region:          "ap-beijing",
				AccessKeyID:     "test-access-key",
				SecretAccessKey: "test-secret-key",
			},
			wantErr: true,
			errMsg:  "endpoint is required for non-AWS S3 services",
		},
		{
			name: "empty type",
			config: &OSSConfig{
				Type:            "",
				Region:          "us-east-1",
				AccessKeyID:     "test-access-key",
				SecretAccessKey: "test-secret-key",
			},
			wantErr: true,
			errMsg:  "storage type is required",
		},
		{
			name: "empty access key",
			config: &OSSConfig{
				Type:            constant.AWSS3,
				Region:          "us-east-1",
				AccessKeyID:     "",
				SecretAccessKey: "test-secret-key",
			},
			wantErr: true,
			errMsg:  "access key ID is required",
		},
		{
			name: "empty secret key",
			config: &OSSConfig{
				Type:            constant.AWSS3,
				Region:          "us-east-1",
				AccessKeyID:     "test-access-key",
				SecretAccessKey: "",
			},
			wantErr: true,
			errMsg:  "secret access key is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOSSConfig_EdgeCases(t *testing.T) {
	t.Run("very long values", func(t *testing.T) {
		longString := string(make([]byte, 1000)) // 1000个空字节
		config := &OSSConfig{
			Type:            constant.AWSS3,
			Region:          longString,
			AccessKeyID:     longString,
			SecretAccessKey: longString,
		}

		err := config.Validate()
		assert.NoError(t, err) // 长字符串应该被接受
	})

	t.Run("special characters in values", func(t *testing.T) {
		config := &OSSConfig{
			Type:            constant.AWSS3,
			Region:          "us-east-1",
			AccessKeyID:     "test@access#key$",
			SecretAccessKey: "test!secret%key^",
		}

		err := config.Validate()
		assert.NoError(t, err) // 特殊字符应该被接受
	})

	t.Run("whitespace only values", func(t *testing.T) {
		config := &OSSConfig{
			Type:            "   ",
			Region:          "us-east-1",
			AccessKeyID:     "test-access-key",
			SecretAccessKey: "test-secret-key",
		}

		err := config.Validate()
		assert.Error(t, err) // 空白字符应该被视为无效
		assert.Contains(t, err.Error(), "storage type is required")
	})
}

func TestOSSConfig_ConfigurationScenarios(t *testing.T) {
	t.Run("production AWS S3 config", func(t *testing.T) {
		config := &OSSConfig{
			Type:            constant.AWSS3,
			Region:          "us-west-2",
			AccessKeyID:     "11",
			SecretAccessKey: "22/33/44",
			MaxRetries:      3,
			ForcePathStyle:  false,
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("development MinIO config", func(t *testing.T) {
		config := &OSSConfig{
			Type:            constant.MinIO,
			Endpoint:        "http://localhost:9000",
			Region:          "us-east-1",
			AccessKeyID:     "minioadmin",
			SecretAccessKey: "minioadmin",
			MaxRetries:      1,
			ForcePathStyle:  true,
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("aliyun OSS config", func(t *testing.T) {
		config := &OSSConfig{
			Type:            constant.AliyunOSS,
			Endpoint:        "https://oss-cn-hangzhou.aliyuncs.com",
			Region:          "cn-hangzhou",
			AccessKeyID:     "11",
			SecretAccessKey: "222",
			MaxRetries:      5,
			ForcePathStyle:  false,
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("tencent COS config", func(t *testing.T) {
		config := &OSSConfig{
			Type:            constant.TencentCOS,
			Endpoint:        "https://cos.ap-beijing.myqcloud.com",
			Region:          "ap-beijing",
			AccessKeyID:     "11",
			SecretAccessKey: "222",
			MaxRetries:      3,
			ForcePathStyle:  false,
		}

		err := config.Validate()
		assert.NoError(t, err)
	})
}
