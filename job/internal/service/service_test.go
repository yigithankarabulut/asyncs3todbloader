package service_test

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/yigithankarabulut/asyncs3todbloader/job/config"
	"github.com/yigithankarabulut/asyncs3todbloader/job/internal/customerror"
	"github.com/yigithankarabulut/asyncs3todbloader/job/internal/service"
	"github.com/yigithankarabulut/asyncs3todbloader/job/internal/storage/objectinfostorage"
	"github.com/yigithankarabulut/asyncs3todbloader/job/internal/storage/productstorage"
	"github.com/yigithankarabulut/asyncs3todbloader/job/model"
	"io"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestService_CheckIfBucketExists(t *testing.T) {
	type fields struct {
		s3Data      config.S3
		s3Client    service.S3Client
		logger      slog.Logger
		s3OutChan   chan *s3.GetObjectOutput
		lineChan    chan string
		productChan chan model.Product
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Empty BucketName should return error",
			fields: fields{
				s3Data: config.S3{
					BucketName: "",
					ObjectKey:  "product.jsonl",
				},
				s3Client: &mockS3Client{
					mockHeadBucket: func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
						if *params.Bucket == "" {
							return nil, customerror.ErrBucketNotFound
						}
						return &s3.HeadBucketOutput{}, nil
					}},
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
		{
			name: "Bucket found, should return success",
			fields: fields{
				s3Data: config.S3{
					BucketName: "test",
					ObjectKey:  "",
				},
				s3Client: &mockS3Client{
					mockHeadBucket: func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
						if *params.Bucket == "test" {
							return &s3.HeadBucketOutput{}, nil
						} else {
							return nil, customerror.ErrBucketNotFound
						}
					},
				},
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.New(
				service.WithS3Data(tt.fields.s3Data),
				service.WithS3Client(tt.fields.s3Client),
			)
			err := s.CheckIfBucketExists(tt.args.ctx)
			if errors.As(err, &customerror.ErrBucketNotFound) != tt.wantErr {
				t.Errorf("CheckIfBucketExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_CheckIfObjectExists(t *testing.T) {

	type fields struct {
		s3Data      config.S3
		s3Client    service.S3Client
		logger      slog.Logger
		s3OutChan   chan *s3.GetObjectOutput
		lineChan    chan string
		productChan chan model.Product
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Invalid ObjectKey should return error",
			fields: fields{
				s3Data: config.S3{
					BucketName: "test",
					ObjectKey:  "123",
				},
				s3Client: &mockS3Client{
					mockHeadObject: func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
						if *params.Bucket == "test" && *params.Key == "test" {
							return &s3.HeadObjectOutput{}, nil
						} else {
							return nil, customerror.ErrObjectNotFound
						}
					},
				},
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
		},
		{
			name: "Empty ObjectKey should return success",
			fields: fields{
				s3Data: config.S3{
					BucketName: "test",
					ObjectKey:  "test",
				},
				s3Client: &mockS3Client{
					mockHeadObject: func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
						if *params.Key == "" {
							return nil, customerror.ErrObjectNotFound
						}
						return &s3.HeadObjectOutput{}, nil
					},
				},
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.New(
				service.WithS3Data(tt.fields.s3Data),
				service.WithS3Client(tt.fields.s3Client),
			)
			err := s.CheckIfObjectExists(tt.args.ctx)
			if errors.As(err, &customerror.ErrObjectNotFound) != tt.wantErr {
				t.Errorf("CheckIfObjectExists() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_CheckObjectDuplicateAndCreate(t *testing.T) {
	type fields struct {
		s3Data            config.S3
		s3Client          service.S3Client
		logger            slog.Logger
		s3OutChan         chan *s3.GetObjectOutput
		lineChan          chan string
		productChan       chan model.Product
		objectInfoStorage objectinfostorage.ObjectInfoStorer
	}
	type args struct {
		ctx context.Context
		out *s3.GetObjectOutput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		errType error
	}{
		{
			name: "out field content type is nil should return error",
			fields: fields{
				s3Data: config.S3{
					BucketName: "test",
					ObjectKey:  "test",
				},
				objectInfoStorage: &mockObjectInfoStorage{},
			},
			args: args{
				ctx: context.Background(),
				out: &s3.GetObjectOutput{
					ContentType:   nil,
					ContentLength: new(int64),
					ETag:          new(string),
				},
			},
			wantErr: true,
			errType: customerror.ErrNilObjectFields,
		},
		{
			name: "out field content length is nil should return error",
			fields: fields{
				s3Data: config.S3{
					BucketName: "test",
					ObjectKey:  "test",
				},
				objectInfoStorage: &mockObjectInfoStorage{},
			},
			args: args{
				ctx: context.Background(),
				out: &s3.GetObjectOutput{
					ContentType:   new(string),
					ContentLength: nil,
					ETag:          new(string),
				},
			},
			wantErr: true,
			errType: customerror.ErrNilObjectFields,
		},
		{
			name: "out field etag is nil should return error",
			fields: fields{
				s3Data: config.S3{
					BucketName: "test",
					ObjectKey:  "test",
				},
				objectInfoStorage: &mockObjectInfoStorage{},
			},
			args: args{
				ctx: context.Background(),
				out: &s3.GetObjectOutput{
					ContentType:   new(string),
					ContentLength: new(int64),
					ETag:          nil,
				},
			},
			wantErr: true,
			errType: customerror.ErrNilObjectFields,
		},
		{
			name: "objectInfoStorage create error should return error",
			fields: fields{
				s3Data: config.S3{
					BucketName: "test",
					ObjectKey:  "test",
				},
				objectInfoStorage: &mockObjectInfoStorage{
					createErr: customerror.ErrCreateObjectInfo,
				},
			},
			args: args{
				ctx: context.Background(),
				out: &s3.GetObjectOutput{
					ContentType:   new(string),
					ContentLength: new(int64),
					ETag:          new(string),
				},
			},
			wantErr: true,
			errType: customerror.ErrCreateObjectInfo,
		},
		{
			name: "out is not nil should return success",
			fields: fields{
				s3Data: config.S3{
					BucketName: "test",
					ObjectKey:  "test",
				},
				objectInfoStorage: &mockObjectInfoStorage{
					createErr: nil,
				},
			},
			args: args{
				ctx: context.Background(),
				out: &s3.GetObjectOutput{
					ContentType:   new(string),
					ContentLength: new(int64),
					ETag:          new(string),
				},
			},
			wantErr: false,
			errType: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.New(
				service.WithS3Data(tt.fields.s3Data),
				service.WithObjectInfoStorage(tt.fields.objectInfoStorage),
			)
			err := s.CheckObjectDuplicateAndCreate(tt.args.ctx, tt.args.out)
			if tt.wantErr && errors.Is(err, tt.errType) {
				t.Errorf("CheckObjectDuplicateAndCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GetObjectFromS3(t *testing.T) {
	type fields struct {
		s3Data      config.S3
		s3Client    service.S3Client
		logger      slog.Logger
		s3OutChan   chan *s3.GetObjectOutput
		lineChan    chan string
		productChan chan model.Product
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		errType error
	}{
		{
			name: "Bucket not found should return error",
			fields: fields{
				s3Data: config.S3{
					BucketName: "test",
					ObjectKey:  "test",
				},
				s3Client: &mockS3Client{
					mockHeadBucket: func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
						return nil, customerror.ErrBucketNotFound
					},
					mockHeadObject: func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
						return &s3.HeadObjectOutput{}, nil
					},
					mockGetObject: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
						return &s3.GetObjectOutput{}, nil
					},
				},
				s3OutChan: make(chan *s3.GetObjectOutput, 1),
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
			errType: customerror.ErrBucketNotFound,
		},
		{
			name: "Object not found should return error",
			fields: fields{
				s3Data: config.S3{
					BucketName: "test",
					ObjectKey:  "test",
				},
				s3Client: &mockS3Client{
					mockHeadBucket: func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
						return &s3.HeadBucketOutput{}, nil
					},
					mockHeadObject: func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
						return nil, customerror.ErrObjectNotFound
					},
					mockGetObject: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
						return &s3.GetObjectOutput{}, nil
					},
				},
				s3OutChan: make(chan *s3.GetObjectOutput, 1),
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
			errType: customerror.ErrObjectNotFound,
		},
		{
			name: "Get object failed should return error",
			fields: fields{
				s3Data: config.S3{
					BucketName: "test",
					ObjectKey:  "test",
				},
				s3Client: &mockS3Client{
					mockHeadBucket: func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
						return &s3.HeadBucketOutput{}, nil
					},
					mockHeadObject: func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
						return &s3.HeadObjectOutput{}, nil
					},
					mockGetObject: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
						return nil, customerror.ErrGetObjectFailed
					},
				},
				s3OutChan: make(chan *s3.GetObjectOutput, 1),
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
			errType: customerror.ErrGetObjectFailed,
		},
		{
			name: "Get object success should return success",
			fields: fields{
				s3Data: config.S3{
					BucketName: "test",
					ObjectKey:  "test",
				},
				s3Client: &mockS3Client{
					mockHeadBucket: func(ctx context.Context, params *s3.HeadBucketInput, optFns ...func(*s3.Options)) (*s3.HeadBucketOutput, error) {
						return &s3.HeadBucketOutput{}, nil
					},
					mockHeadObject: func(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
						return &s3.HeadObjectOutput{}, nil
					},
					mockGetObject: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
						return &s3.GetObjectOutput{}, nil
					},
				},
				s3OutChan: make(chan *s3.GetObjectOutput, 1),
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
			errType: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.New(
				service.WithS3Data(tt.fields.s3Data),
				service.WithS3Client(tt.fields.s3Client),
				service.WithS3OutChan(tt.fields.s3OutChan),
			)
			err := s.GetObjectFromS3(tt.args.ctx)
			if tt.wantErr && errors.Is(err, tt.errType) {
				t.Errorf("GetObjectFromS3() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_ReadDataFromS3Object(t *testing.T) {
	type fields struct {
		s3Data      config.S3
		s3Client    service.S3Client
		logger      slog.Logger
		s3OutChan   chan *s3.GetObjectOutput
		lineChan    chan string
		productChan chan model.Product
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		routineFunc func(a ...interface{})
		wantErr     bool
		errType     error
	}{
		{
			name: "s3OutChan is closed should return error",
			fields: fields{
				s3OutChan: make(chan *s3.GetObjectOutput, 1),
				lineChan:  make(chan string, 10),
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
			errType: customerror.ErrChannelClosed,
			routineFunc: func(a ...interface{}) {
				close(a[0].(chan *s3.GetObjectOutput))
			},
		},
		{
			name: "File scan failed should return error",
			fields: fields{
				s3OutChan: make(chan *s3.GetObjectOutput, 1),
				lineChan:  make(chan string, 10),
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
			errType: customerror.ErrFileScanFailed,
			routineFunc: func(a ...interface{}) {
				a[0].(chan *s3.GetObjectOutput) <- &s3.GetObjectOutput{
					Body: io.NopCloser(strings.NewReader("")),
				}
			},
		},
		{
			name: "Success should return success",
			fields: fields{
				s3OutChan: make(chan *s3.GetObjectOutput, 1),
				lineChan:  make(chan string, 10),
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
			errType: nil,
			routineFunc: func(a ...interface{}) {
				a[0].(chan *s3.GetObjectOutput) <- &s3.GetObjectOutput{
					Body: io.NopCloser(strings.NewReader("test")),
				}
				close(a[0].(chan *s3.GetObjectOutput))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.New(
				service.WithS3OutChan(tt.fields.s3OutChan),
				service.WithLineChannel(tt.fields.lineChan),
				service.WithLogger(slog.New(
					slog.NewJSONHandler(os.Stdout, nil),
				)),
			)
			go func() {
				tt.routineFunc(tt.fields.s3OutChan)
			}()
			err := s.ReadDataFromS3Object(tt.args.ctx)
			if tt.wantErr && errors.Is(err, tt.errType) {
				t.Errorf("ReadDataFromS3Object() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_HandleLines(t *testing.T) {
	type fields struct {
		s3Data      config.S3
		s3Client    service.S3Client
		logger      slog.Logger
		lineChan    chan string
		productChan chan model.Product
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		routineFunc func(a ...interface{})
		wantErr     bool
		errType     error
	}{
		{
			name: "lineChan is closed should return error",
			fields: fields{
				lineChan:    make(chan string, 10),
				productChan: make(chan model.Product, 10),
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
			errType: customerror.ErrChannelClosed,
			routineFunc: func(a ...interface{}) {
				close(a[0].(chan string))
			},
		},
		{
			name: "json unmarshal failed but should continue and error log should be printed",
			fields: fields{
				lineChan:    make(chan string, 10),
				productChan: make(chan model.Product, 10),
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
			errType: nil,
			routineFunc: func(a ...interface{}) {
				a[0].(chan string) <- "test"
				close(a[0].(chan string))
			},
		},
		{
			name: "Success should return success",
			fields: fields{
				lineChan:    make(chan string, 10),
				productChan: make(chan model.Product, 10),
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
			errType: nil,
			routineFunc: func(a ...interface{}) {
				a[0].(chan string) <- `{"id":1}`
				close(a[0].(chan string))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.New(
				service.WithLineChannel(tt.fields.lineChan),
				service.WithProductChannel(tt.fields.productChan),
				service.WithLogger(slog.New(
					slog.NewJSONHandler(os.Stdout, nil),
				)),
				service.WithLineHandlerWorkerCount(1),
			)
			go func() {
				tt.routineFunc(tt.fields.lineChan, tt.fields.productChan)
			}()
			err := s.HandleLines(tt.args.ctx)
			if tt.wantErr && errors.Is(err, tt.errType) {
				t.Errorf("HandleLines() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_WriteDataToDb(t *testing.T) {
	type fields struct {
		s3Data         config.S3
		s3Client       service.S3Client
		logger         slog.Logger
		productStorage productstorage.ProductStorer
		productChan    chan model.Product
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		routineFunc func(a ...interface{})
		wantErr     bool
		errType     error
	}{
		{
			name: "productChan is closed should return error",
			fields: fields{
				productChan: make(chan model.Product, 10),
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
			errType: customerror.ErrChannelClosed,
			routineFunc: func(a ...interface{}) {
				close(a[0].(chan model.Product))
			},
		},
		{
			name: "context is done should return error",
			fields: fields{
				productChan: make(chan model.Product, 10),
			},
			args:    args{ctx: context.Background()},
			wantErr: true,
			errType: context.Canceled,
			routineFunc: func(a ...interface{}) {
				close(a[0].(chan model.Product))
			},
		},
		{
			name: "Create product failed but should continue and error log should be printed",
			fields: fields{
				productChan: make(chan model.Product, 10),
				productStorage: &mockProductStorage{
					createErr: customerror.ErrCreateProduct,
				},
			},
			args:    args{ctx: context.Background()},
			wantErr: false,
			errType: nil,
			routineFunc: func(a ...interface{}) {
				a[0].(chan model.Product) <- model.Product{}
				close(a[0].(chan model.Product))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.New(
				service.WithProductChannel(tt.fields.productChan),
				service.WithDBWriteWorkerCount(1),
				service.WithProductStorage(tt.fields.productStorage),
				service.WithLogger(slog.New(
					slog.NewJSONHandler(os.Stdout, nil),
				)),
			)
			go func() {
				tt.routineFunc(tt.fields.productChan)
			}()
			err := s.WriteDataToDb(tt.args.ctx)
			if tt.wantErr && errors.Is(err, tt.errType) {
				t.Errorf("WriteDataToDb() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
