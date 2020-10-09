package server

import (
	"log"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/user-service/pkg/middlewares"
	. "github.com/user-service/pkg/models"
	. "github.com/user-service/pkg/rpc"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type UserService interface {
	StartAPI()
	StartRPC()
	waitShutdown(chan error, chan os.Signal)
}

type Service struct {
	controllers  fasthttp.RequestHandler
	port         string
	rpcPort      string
	interruptAPI chan os.Signal
	listenAPI    chan error
	interruptRPC chan os.Signal
	listenRPC    chan error
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
		controllers:  initControllers().Handler,
		port:         os.Getenv("PORT"),
		rpcPort:      os.Getenv("RPC_PORT"),
		interruptAPI: make(chan os.Signal, 1),
		listenAPI:    make(chan error, 1),
		interruptRPC: make(chan os.Signal, 1),
		listenRPC:    make(chan error, 1),
	}
}

func (srv *Service) StartRPC() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	api := new(UserRPC)
	err := rpc.Register(api)
	if err != nil {
		zap.S().Fatalf("Listener error: %v", err)
	}
	rpc.HandleHTTP()

	zap.S().Info("RPC started on port: " + srv.rpcPort)

	go func(listenRPC chan error) {
		zap.S().Info("Service started on port: " + srv.port)

		listenRPC <- http.ListenAndServe(":"+srv.rpcPort, nil)
	}(srv.listenRPC)

	signal.Notify(srv.interruptAPI, syscall.SIGINT, syscall.SIGTERM)
	srv.waitShutdown(srv.listenRPC, srv.interruptRPC)
}

func (srv *Service) StartAPI() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go func(listenAPI chan error) {
		zap.S().Info("Service started on port: " + srv.port)
		listenAPI <- fasthttp.ListenAndServe(":"+srv.port, middlewares.CORS(srv.controllers))
	}(srv.listenAPI)

	signal.Notify(srv.interruptAPI, syscall.SIGINT, syscall.SIGTERM)
	srv.waitShutdown(srv.listenAPI, srv.interruptAPI)
}

func (srv Service) waitShutdown(listen chan error, interrupt chan os.Signal) {
	for {
		select {
		case err := <-listen:
			if err != nil {
				zap.S().Fatalf("Listener error: %v", err)
			}
			os.Exit(0)
		case err := <-interrupt:
			zap.S().Fatalf("Shutdown signal: %v", err.String())
		}
	}
}
