package app

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/acool-kaz/forum-auth-service/internal/config"
	grpcAuthHandler "github.com/acool-kaz/forum-auth-service/internal/delivery/grpc"
	"github.com/acool-kaz/forum-auth-service/internal/repository"
	"github.com/acool-kaz/forum-auth-service/internal/service"
	auth_svc_pb "github.com/acool-kaz/forum-auth-service/pkg/auth_svc/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type app struct {
	cfg *config.Config

	db *sql.DB

	grpcServer      *grpc.Server
	grpcAuthHandler *grpcAuthHandler.AuthSvcHandler
}

func InitApp(cfg *config.Config) (*app, error) {
	log.Println("init app")
	db, err := repository.InitDB(cfg)
	if err != nil {
		return nil, fmt.Errorf("init app: %w", err)
	}

	repo := repository.InitRepository(db)
	service := service.InitService(repo)

	grpcAuthHandler := grpcAuthHandler.InitAuthSvcHandler(service)

	return &app{
		cfg:             cfg,
		db:              db,
		grpcAuthHandler: grpcAuthHandler,
	}, nil
}

func (a *app) RunApp() {
	log.Println("run app")

	go func() {
		if err := a.startGRPC(); err != nil {
			log.Println(err)
			return
		}
	}()
	log.Println("grpc started on", a.cfg.Grpc.Port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	fmt.Println()
	log.Println("Received terminate, graceful shutdown", sig)

	log.Println("grpc: Server closed")

	if err := a.db.Close(); err != nil {
		log.Println(err)
	} else {
		log.Println("db closed")
	}
}

func (a *app) startGRPC() error {
	listen, err := net.Listen(a.cfg.Grpc.Type, a.cfg.Grpc.Host+":"+a.cfg.Grpc.Port)
	if err != nil {
		return err
	}
	defer listen.Close()

	opt := []grpc.ServerOption{}

	a.grpcServer = grpc.NewServer(opt...)

	auth_svc_pb.RegisterAuthServiceServer(a.grpcServer, a.grpcAuthHandler)

	reflection.Register(a.grpcServer)

	return a.grpcServer.Serve(listen)
}
