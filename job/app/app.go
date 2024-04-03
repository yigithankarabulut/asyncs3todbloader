package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	appConfig "github.com/yigithankarabulut/asyncs3todbloader/job/config"
	"github.com/yigithankarabulut/asyncs3todbloader/job/internal/customerror"
	"github.com/yigithankarabulut/asyncs3todbloader/job/internal/service"
	"github.com/yigithankarabulut/asyncs3todbloader/job/internal/storage/objectinfostorage"
	"github.com/yigithankarabulut/asyncs3todbloader/job/internal/storage/productstorage"
	"github.com/yigithankarabulut/asyncs3todbloader/job/model"
	"github.com/yigithankarabulut/asyncs3todbloader/job/pkg/mongo"
	"log/slog"
	"os"
	"sync"
	"time"
)

const (
	LineChannelSize    = 50
	ProductChannelSize = 50
	LineHandlerCount   = 50
	DBWriteWorkerCount = 50
)

type app struct {
	config   *appConfig.Config
	logLevel slog.Level
	logger   *slog.Logger
	doneChan chan struct{}
}

type Option func(*app)

func WithLogLevel(level string) Option {
	return func(s *app) {
		switch level {
		case "DEBUG":
			s.logLevel = slog.LevelDebug
		case "INFO":
			s.logLevel = slog.LevelInfo
		case "WARN":
			s.logLevel = slog.LevelWarn
		case "ERROR":
			s.logLevel = slog.LevelError
		default:
			s.logLevel = slog.LevelInfo
		}
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(s *app) {
		s.logger = logger
	}
}

func WithConfig(cfg *appConfig.Config) Option {
	return func(s *app) {
		s.config = cfg
	}
}

func WithDoneChan(doneChan chan struct{}) Option {
	return func(s *app) {
		s.doneChan = doneChan
	}
}

// New creates a new app instance. It initializes the storages, logger, connects to MongoDB and AWS, and runs the app.
func New(opts ...Option) error {
	app := &app{
		logLevel: slog.LevelInfo,
	}
	for _, opt := range opts {
		opt(app)
	}
	if app.config == nil {
		return errors.New("config is required")
	}
	// set default logger if not provided
	if app.logger == nil {
		logHandlerOpts := &slog.HandlerOptions{Level: app.logLevel}
		logHandler := slog.NewJSONHandler(os.Stdout, logHandlerOpts)
		app.logger = slog.New(logHandler)
	}
	slog.SetDefault(app.logger)
	app.logger.Info("Starting app...")

	// Connect to MongoDB
	db, err := mongo.ConnectMongo(app.config.Database)
	if err != nil {
		return fmt.Errorf("error connecting to mongo: %w", err)
	}

	// Connect to AWS
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	awsConfig, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(app.config.Aws.Region),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     app.config.Aws.AccessKey,
				SecretAccessKey: app.config.Aws.SecretKey,
			}}),
	)
	if err != nil {
		return fmt.Errorf("error loading aws config: %w", err)
	}
	s3Client := s3.NewFromConfig(awsConfig)

	// Initialize storage instances
	productStorage := productstorage.New(
		productstorage.WithProductCollection(app.config.Database.ProductCollection),
		productstorage.WithDB(db),
	)
	objectInfoStorage := objectinfostorage.New(
		objectinfostorage.WithObjectCollection(app.config.Database.ObjectInfoCollection),
		objectinfostorage.WithDB(db),
	)

	// Create indexes
	if err := productStorage.CreateIndex(context.Background()); err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	if err := objectInfoStorage.CreateIndex(context.Background()); err != nil {
		return fmt.Errorf("error creating index: %w", err)
	}
	return app.Run(s3Client, productStorage, objectInfoStorage)
}

// Run starts the app. It processes each S3 object concurrently.
// It creates a service instance for each S3 object and runs it.
// Each service instance will have its own out, line, and product channels.
// Each S3 object will have its own line handler and db writer workers.
// It waits for all S3 objects to be processed and sends a signal to the done channel.
func (a *app) Run(s3Client *s3.Client, productStorage productstorage.ProductStorer, objectInfoStorage objectinfostorage.ObjectInfoStorer) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	start := time.Now()

	wg := sync.WaitGroup{}
	for _, s3Object := range a.config.Aws.S3 {
		wg.Add(1)
		go func(s3Object appConfig.S3) {
			defer wg.Done()
			outChan := make(chan *s3.GetObjectOutput, 1)
			lineChan := make(chan string, LineChannelSize)
			productChan := make(chan model.Product, ProductChannelSize)

			service := service.New(
				service.WithS3Client(s3Client),
				service.WithS3Data(s3Object),
				service.WithProductStorage(productStorage),
				service.WithObjectInfoStorage(objectInfoStorage),
				service.WithLogger(a.logger),
				service.WithS3OutChan(outChan),
				service.WithProductChannel(productChan),
				service.WithLineChannel(lineChan),
				service.WithLineHandlerWorkerCount(LineHandlerCount),
				service.WithDBWriteWorkerCount(DBWriteWorkerCount),
			)

			if err := service.Run(); err != nil {
				var ce *customerror.Error
				// check if error is custom error.
				// if error must be logged, log it.
				if errors.As(err, &ce) {
					message := ce.Message
					if ce.Data != nil {
						data, ok := ce.Data.(string)
						if ok {
							message += ", " + data
						}
						if ce.Loggable {
							a.logger.Error(message)
						}
					}
				}
				return
			}
		}(s3Object)
	}
	wg.Wait()
	a.doneChan <- struct{}{}
	elapsed := time.Since(start)
	a.logger.Info(fmt.Sprintf("Elapsed Time: %s", elapsed))
	return nil
}
