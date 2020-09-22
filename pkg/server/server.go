package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/user-service/pkg/middlewares"
	. "github.com/user-service/pkg/models"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type UserService interface {
	StartService()
	waitShutdown()
}

type Service struct {
	controllers fasthttp.RequestHandler
	port        string
	interrupt   chan os.Signal
	listen      chan error
}

func init() {
	err := godotenv.Load()
	if err != nil {
		zap.S().Fatalf("Load env error: %v", err)
	}

	err = initLogger()
	if err != nil {
		zap.S().Fatalf("InitLogger error: %v", err)
	}

	err = InitDatabase()
	if err != nil {
		zap.S().Fatalf("InitDatabase error: %v", err)
	}
	zap.S().Info("Database connect to " + os.Getenv("DB_NAME"))

	err = InitRedis()
	if err != nil {
		zap.S().Fatalf("InitRedis error: %v", err)
	}
	zap.S().Info("Redis connect to " + os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"))

	port := os.Getenv("PORT")
	if port == "" {
		zap.S().Fatalf("PORT env does not exist")
	}
}

func NewService() *Service {
	return &Service{
		controllers: initControllers().Handler,
		port:        os.Getenv("PORT"),
		interrupt:   make(chan os.Signal, 1),
		listen:      make(chan error, 1),
	}
}

func (srv *Service) StartService() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go func(listen chan error) {
		zap.S().Info("Service started on port: " + srv.port)
		listen <- fasthttp.ListenAndServe(":"+srv.port, middlewares.CORS(srv.controllers))
	}(srv.listen)

	signal.Notify(srv.interrupt, syscall.SIGINT, syscall.SIGTERM)
	srv.waitShutdown()
}

func (srv *Service) waitShutdown() {
	for {
		select {
		case err := <-srv.listen:
			if err != nil {
				zap.S().Fatalf("Listener error: %v", err)
			}
			os.Exit(0)
		case err := <-srv.interrupt:
			zap.S().Fatalf("Shutdown signal: %v", err.String())
		}
	}
}
