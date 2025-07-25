package initialize

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/supuwoerc/weaver/conf"
	"github.com/supuwoerc/weaver/pkg/constant"
)

// MockOSSClient 用于单元测试的mock S3客户端
// 只实现部分接口方法用于测试
type MockOSSClient struct {
	headBucketCalled bool
}

func (m *MockOSSClient) HeadBucket(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
	m.headBucketCalled = true
	return &s3.HeadBucketOutput{}, nil
}

// 其它接口方法可根据需要补充
func (m *MockOSSClient) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	return &s3.HeadObjectOutput{}, nil
}
func (m *MockOSSClient) ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error) {
	return &s3.ListBucketsOutput{}, nil
}
func (m *MockOSSClient) ListDirectoryBuckets(ctx context.Context, params *s3.ListDirectoryBucketsInput, optFns ...func(*s3.Options)) (*s3.ListDirectoryBucketsOutput, error) {
	return &s3.ListDirectoryBucketsOutput{}, nil
}
func (m *MockOSSClient) ListMultipartUploads(ctx context.Context, params *s3.ListMultipartUploadsInput, optFns ...func(*s3.Options)) (*s3.ListMultipartUploadsOutput, error) {
	return &s3.ListMultipartUploadsOutput{}, nil
}
func (m *MockOSSClient) ListObjectVersions(ctx context.Context, params *s3.ListObjectVersionsInput, optFns ...func(*s3.Options)) (*s3.ListObjectVersionsOutput, error) {
	return &s3.ListObjectVersionsOutput{}, nil
}
func (m *MockOSSClient) ListParts(ctx context.Context, params *s3.ListPartsInput, optFns ...func(*s3.Options)) (*s3.ListPartsOutput, error) {
	return &s3.ListPartsOutput{}, nil
}

func TestNewS3CompatibleStorage(t *testing.T) {
	cfg := &conf.OSSConfig{
		Type:            constant.AWSS3,
		Region:          "us-east-1",
		AccessKeyID:     "test-access-key",
		SecretAccessKey: "test-secret-key",
	}
	mockClient := &MockOSSClient{}
	storage := NewS3CompatibleStorage(cfg, mockClient)
	assert.NotNil(t, storage)
	assert.Equal(t, cfg, storage.GetConfig())
	assert.Equal(t, mockClient, storage.GetClient())
}

func TestS3CompatibleStorage_GetClientAndConfig(t *testing.T) {
	cfg := &conf.OSSConfig{
		Type:            constant.AWSS3,
		Region:          "us-east-1",
		AccessKeyID:     "test-access-key",
		SecretAccessKey: "test-secret-key",
	}
	mockClient := &MockOSSClient{}
	storage := NewS3CompatibleStorage(cfg, mockClient)
	assert.Equal(t, mockClient, storage.GetClient())
	assert.Equal(t, cfg, storage.GetConfig())
}

func TestNewS3Client_AwsS3(t *testing.T) {
	cfg := &conf.OSSConfig{
		Type:            constant.AWSS3,
		Region:          "us-east-1",
		AccessKeyID:     "test-access-key",
		SecretAccessKey: "test-secret-key",
	}
	client := NewS3Client(cfg)
	assert.NotNil(t, client)
}

func TestNewS3Client_MinIO(t *testing.T) {
	cfg := &conf.OSSConfig{
		Type:            constant.MinIO,
		Endpoint:        "http://localhost:9000",
		Region:          "us-east-1",
		AccessKeyID:     "minioadmin",
		SecretAccessKey: "minioadmin",
	}
	client := NewS3Client(cfg)
	assert.NotNil(t, client)
}

func TestNewS3Client_InvalidType(t *testing.T) {
	cfg := &conf.OSSConfig{
		Type:            "invalid-type",
		Region:          "us-east-1",
		AccessKeyID:     "test-access-key",
		SecretAccessKey: "test-secret-key",
	}
	assert.Panics(t, func() {
		NewS3Client(cfg)
	})
}

func TestNewS3Client_InvalidConfig(t *testing.T) {
	cfg := &conf.OSSConfig{
		Type:            constant.AWSS3,
		Region:          "us-east-1",
		AccessKeyID:     "",
		SecretAccessKey: "",
	}
	assert.Panics(t, func() {
		NewS3Client(cfg)
	})
}
