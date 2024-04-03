package service

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/yigithankarabulut/asyncs3todbloader/job/config"
	"github.com/yigithankarabulut/asyncs3todbloader/job/internal/storage/objectinfostorage"
	"github.com/yigithankarabulut/asyncs3todbloader/job/internal/storage/productstorage"
	"github.com/yigithankarabulut/asyncs3todbloader/job/model"
	"log/slog"
)

type Service interface {
	Run() error
	CheckIfBucketExists(ctx context.Context) error
	CheckIfObjectExists(ctx context.Context) error
	CheckObjectDuplicateAndCreate(ctx context.Context, out *s3.GetObjectOutput) error
	GetObjectFromS3(ctx context.Context) error
	ReadDataFromS3Object(ctx context.Context) error
	HandleLines(ctx context.Context) error
	WriteDataToDb(ctx context.Context) error
}

type S3Client interface {
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	HeadBucket(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type service struct {
	s3Data                 config.S3
	s3Client               S3Client
	logger                 *slog.Logger
	productStorage         productstorage.ProductStorer
	objectInfoStorage      objectinfostorage.ObjectInfoStorer
	s3OutChan              chan *s3.GetObjectOutput
	lineChan               chan string
	productChan            chan model.Product
	lineHandlerWorkerCount int
	dbWriteWorkerCount     int
}

type Option func(*service)

func WithS3Client(s3Client S3Client) Option {
	return func(s *service) {
		s.s3Client = s3Client
	}
}

func WithS3Data(s3Data config.S3) Option {
	return func(s *service) {
		s.s3Data = s3Data
	}
}

func WithProductStorage(productStorage productstorage.ProductStorer) Option {
	return func(s *service) {
		s.productStorage = productStorage
	}
}

func WithObjectInfoStorage(objectInfoStorage objectinfostorage.ObjectInfoStorer) Option {
	return func(s *service) {
		s.objectInfoStorage = objectInfoStorage
	}
}

func WithS3OutChan(ch chan *s3.GetObjectOutput) Option {
	return func(s *service) {
		s.s3OutChan = ch
	}
}

func WithProductChannel(ch chan model.Product) Option {
	return func(s *service) {
		s.productChan = ch
	}
}

func WithLineChannel(ch chan string) Option {
	return func(s *service) {
		s.lineChan = ch
	}
}

func WithLineHandlerWorkerCount(count int) Option {
	return func(s *service) {
		s.lineHandlerWorkerCount = count
	}
}

func WithDBWriteWorkerCount(count int) Option {
	return func(s *service) {
		s.dbWriteWorkerCount = count
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(s *service) {
		s.logger = logger
	}
}

func New(opts ...Option) Service {
	s := &service{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}
