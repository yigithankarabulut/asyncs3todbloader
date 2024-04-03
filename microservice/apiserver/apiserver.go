package apiserver

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/config"
	productrepository "github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/repository/product"
	productservice "github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/service/product"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/transport/http/basehttphandler"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/internal/transport/http/producthandler"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/pkg"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/pkg/middleware"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/pkg/mongo"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/pkg/response"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/pkg/validator"
	"github.com/yigithankarabulut/asyncs3todbloader/microservice/releaseinfo"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ContextCancelTimeout = 5 * time.Second
	ShutdownTimeout      = 5 * time.Second
	ServerReadTimeout    = 5 * time.Second
	ServerWriteTimeout   = 5 * time.Second
	ServerIdleTimeout    = 5 * time.Second
)

type HttpEndpoints interface {
	AddRoutes(router fiber.Router)
}

type apiServer struct {
	config    *config.Config
	app       *fiber.App
	handlers  []HttpEndpoints
	logLevel  slog.Level
	logger    *slog.Logger
	serverEnv string
}

type Option func(*apiServer)

// WithLogLevel sets the log level option.
func WithLogLevel(level string) Option {
	return func(s *apiServer) {
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

// WithLogger sets the logger option.
func WithLogger(logger *slog.Logger) Option {
	return func(s *apiServer) {
		s.logger = logger
	}
}

// WithServerEnv sets the server environment option.
func WithServerEnv(env string) Option {
	return func(s *apiServer) {
		s.serverEnv = env
	}
}

// WithConfig sets the config option.
func WithConfig(config *config.Config) Option {
	return func(s *apiServer) {
		s.config = config
	}
}

// healthzCheck is a health check endpoint.
func (s *apiServer) healthzCheck() {
	s.app.Get("/healthz/live", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"server":            s.serverEnv,
			"version":           releaseinfo.Version,
			"build_information": releaseinfo.BuildInformation,
			"message":           "liveness is OK!, server is ready to accept connections",
		})
	})
	s.app.Get("/healthz/ping", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "pong",
		})
	})
}

// AppendLayers appends the layers to the server.
func (s *apiServer) AppendLayers() {
	packages := pkg.New(
		pkg.WithValidator(validator.New()),
		pkg.WithResponse(response.New()),
	)
	productRepository := productrepository.New(
		productrepository.WithDB(mongo.DB),
		productrepository.WithProductCollection(s.config.Database.Collection.Name),
	)
	productService := productservice.New(
		productservice.WithRepository(productRepository),
	)
	baseHttpHandler := basehttphandler.New(
		basehttphandler.WithPackages(packages),
		basehttphandler.WithLogger(s.logger),
		basehttphandler.WithContextTimeout(ContextCancelTimeout),
	)
	productHandler := producthandler.New(
		producthandler.WithBaseHttpHandler(baseHttpHandler),
		producthandler.WithProductService(productService),
	)

	s.handlers = append(s.handlers, productHandler)
	for _, handler := range s.handlers {
		handler.AddRoutes(s.app)
	}
}

// New creates a new api server.
func New(opts ...Option) error {
	apiserv := &apiServer{
		logLevel: slog.LevelInfo,
	}
	for _, opt := range opts {
		opt(apiserv)
	}
	if apiserv.config == nil {
		return fmt.Errorf("config is required")
	}
	if err := mongo.Connect(apiserv.config.Database); err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	if apiserv.logger == nil {
		logHandlerOpts := &slog.HandlerOptions{Level: apiserv.logLevel}
		logHandler := slog.NewJSONHandler(os.Stdout, logHandlerOpts)
		apiserv.logger = slog.New(logHandler)
	}
	slog.SetDefault(apiserv.logger)
	if apiserv.serverEnv == "" {
		apiserv.serverEnv = "development"
	}
	apiserv.app = fiber.New(fiber.Config{
		ReadTimeout:  ServerReadTimeout,
		WriteTimeout: ServerWriteTimeout,
		IdleTimeout:  ServerIdleTimeout,
	})
	apiserv.app.Use(recover.New())
	apiserv.app.Use(middleware.HttpLoggingMiddleware(apiserv.logger, apiserv.app))
	apiserv.healthzCheck()
	apiserv.AppendLayers()
	return apiserv.ListenAndServe()
}

// ListenAndServe starts the fiber app and listens for incoming requests. It also listens for shutdown signals and handles graceful shutdown.
func (s *apiServer) ListenAndServe() error {
	shutdown := make(chan os.Signal, 1)
	apiErr := make(chan error, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s.logger.Info("starting api server", "listening on", s.config.Port, "env", s.serverEnv)
		apiErr <- s.app.Listen(":" + s.config.Port)
	}()

	select {
	case err := <-apiErr:
		return fmt.Errorf("error listening api server: %w", err)
	case <-shutdown:
		s.logger.Info("starting shutdown", "pid", os.Getpid())
		ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer cancel()
		if err := s.app.ShutdownWithContext(ctx); err != nil {
			return fmt.Errorf("error shutting down server: %w", err)
		}
		s.logger.Info("shutdown complete", "pid", os.Getpid())
	}
	return nil
}
