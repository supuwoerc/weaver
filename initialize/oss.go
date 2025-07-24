package initialize

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/pkg/constant"
)

type S3CompatibleStorage struct {
	client *s3.Client
	config *conf.OSSConfig
}

// NewS3CompatibleStorage 创建兼容S3协议的存储客户端
func NewS3CompatibleStorage(cfg *conf.OSSConfig) *S3CompatibleStorage {
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
	client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		// 对于MinIO等需要强制路径样式
		if cfg.ForcePathStyle {
			o.UsePathStyle = true
		}
	})
	return &S3CompatibleStorage{
		client: client,
		config: cfg,
	}
}

// GetClient 获取S3客户端
func (s *S3CompatibleStorage) GetClient() *s3.Client {
	return s.client
}

// GetConfig 获取配置
func (s *S3CompatibleStorage) GetConfig() *conf.OSSConfig {
	return s.config
}

// Ping 测试连接
func (s *S3CompatibleStorage) Ping() error {
	// 设置超时上下文
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := s.client.HeadBucket(timeoutCtx, &s3.HeadBucketInput{
		Bucket: aws.String("default"),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to default bucket %v", err)
	}
	return nil
}
