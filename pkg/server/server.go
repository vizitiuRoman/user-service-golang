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
	Init()
	waitShutdown(chan error, chan os.Signal)
}

type Service struct {
	controllers  fasthttp.RequestHandler
	logger       *zap.SugaredLogger
	port         string
	rpcPort      string
	interruptAPI chan os.Signal
	listenAPI    chan error
	interruptRPC chan os.Signal
	listenRPC    chan error
}

func (srv *Service) Init() {
	err := godotenv.Load()
	if err != nil {
		srv.logger.Fatalf("Load env error: %v", err)
	}

	err = InitDatabase()
	if err != nil {
		srv.logger.Fatalf("InitDatabase error: %v", err)
	}
	srv.logger.Info("Database connect to " + os.Getenv("DB_NAME"))

	err = InitRedis()
	if err != nil {
		srv.logger.Fatalf("InitRedis error: %v", err)
	}
	srv.logger.Info("Redis connect to " + os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"))

	port := os.Getenv("PORT")
	if port == "" {
		srv.logger.Fatalf("PORT env does not exist")
	}

	rpcPort := os.Getenv("RPC_PORT")
	if rpcPort == "" {
		srv.logger.Fatalf("RPC_PORT env does not exist")
	}

	srv.port = os.Getenv("PORT")
	srv.rpcPort = os.Getenv("RPC_PORT")
}

func NewService() *Service {
	return &Service{
		controllers:  initControllers().Handler,
		logger:       initLogger(),
		port:         "",
		rpcPort:      "",
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
		srv.logger.Fatalf("Listener error: %v", err)
	}
	rpc.HandleHTTP()

	go func(listenRPC chan error) {
		srv.logger.Info("RPC started on port: " + srv.rpcPort)

		listenRPC <- http.ListenAndServe(":"+srv.rpcPort, nil)
	}(srv.listenRPC)

	signal.Notify(srv.interruptAPI, syscall.SIGINT, syscall.SIGTERM)
	srv.waitShutdown(srv.listenRPC, srv.interruptRPC)
}

func (srv *Service) StartAPI() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	go func(listenAPI chan error) {
		srv.logger.Info("Service started on port: " + srv.port)
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
				srv.logger.Fatalf("Listener error: %v", err)
			}
			os.Exit(0)
		case err := <-interrupt:
			srv.logger.Fatalf("Shutdown signal: %v", err.String())
		}
	}
}
