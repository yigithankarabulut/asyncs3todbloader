package service_test

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/yigithankarabulut/asyncs3todbloader/job/model"
)

var (
	errProductStorageCreate         = errors.New("product storage create error")
	errProductStorageCreateIndex    = errors.New("product storage create index error")
	errObjectInfoStorageCreate      = errors.New("object info storage create error")
	errObjectInfoStorageCreateIndex = errors.New("object info storage create index error")
)

type mockProductStorage struct {
	createIndexErr error
	createErr      error
	createBatchErr error
}

func (m *mockProductStorage) CreateIndex(ctx context.Context) error {
	return m.createIndexErr
}

func (m *mockProductStorage) Create(ctx context.Context, product model.Product) error {
	return m.createErr
}

type mockObjectInfoStorage struct {
	createIndexErr error
	createErr      error
}

func (m *mockObjectInfoStorage) CreateIndex(ctx context.Context) error {
	return m.createIndexErr
}

func (m *mockObjectInfoStorage) Create(ctx context.Context, objectPartition model.ObjectInfo) error {
	return m.createErr
}

type mockS3Client struct {
	mockHeadBucket func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
	mockHeadObject func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	mockGetObject  func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func (m *mockS3Client) HeadBucket(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
	return m.mockHeadBucket(ctx, params)
}

func (m *mockS3Client) HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	return m.mockHeadObject(ctx, params)
}

func (m *mockS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return m.mockGetObject(ctx, params)
}
