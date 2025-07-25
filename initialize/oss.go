package initialize

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/pkg/constant"
)

// OSSClient 定义兼容S3协议的客户端接口，便于mock和单元测试
type OSSClient interface {
	s3.HeadBucketAPIClient
	s3.HeadObjectAPIClient
	s3.ListBucketsAPIClient
	s3.ListDirectoryBucketsAPIClient
	s3.ListMultipartUploadsAPIClient
	s3.ListObjectVersionsAPIClient
	s3.ListPartsAPIClient
}

type S3CompatibleStorage struct {
	client OSSClient
	config *conf.OSSConfig
}

// NewS3CompatibleStorage 创建兼容S3协议的存储客户端
func NewS3CompatibleStorage(cfg *conf.OSSConfig, client OSSClient) *S3CompatibleStorage {
	return &S3CompatibleStorage{
		client: client,
		config: cfg,
	}
}

// GetClient 获取S3客户端
func (s *S3CompatibleStorage) GetClient() OSSClient {
	return s.client
}

// GetConfig 获取配置
func (s *S3CompatibleStorage) GetConfig() *conf.OSSConfig {
	return s.config
}

// NewS3Client 创建兼容s3的客户端
func NewS3Client(cfg *conf.OSSConfig) *s3.Client {
	// 验证配置
	if err := cfg.Validate(); err != nil {
		panic(err)
	}
	var awsConfig aws.Config
	var err error
	// 创建自定义凭证提供者
	credentialsProvider := credentials.NewStaticCredentialsProvider(
		cfg.AccessKeyID,
		cfg.SecretAccessKey,
		"",
	)
	// 设置默认值
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 3
	}
	switch cfg.Type {
	case constant.AWSS3:
		// AWS S3 - 使用默认配置
		awsConfig, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(cfg.Region),
			config.WithCredentialsProvider(credentialsProvider),
			config.WithRetryMode(aws.RetryModeAdaptive),
			config.WithRetryMaxAttempts(cfg.MaxRetries),
		)

	case constant.AliyunOSS, constant.TencentCOS, constant.MinIO:
		// 其他S3兼容服务-使用自定义endpoint
		endpoint := cfg.Endpoint
		awsConfig, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(cfg.Region),
			config.WithBaseEndpoint(endpoint),
			config.WithCredentialsProvider(credentialsProvider),
			config.WithRetryMode(aws.RetryModeAdaptive),
			config.WithRetryMaxAttempts(cfg.MaxRetries),
		)

	default:
		panic(fmt.Errorf("unsupported storage type: %s", cfg.Type))
	}
	if err != nil {
		panic(fmt.Errorf("failed to load AWS config: %w", err))
	}
	// 创建S3客户端
	return s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		// 对于MinIO等需要强制路径样式
		if cfg.ForcePathStyle {
			o.UsePathStyle = true
		}
	})
}
