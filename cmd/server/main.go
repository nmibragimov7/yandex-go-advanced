package server

import (
	"log"
	"net"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/grpc_handlers"
	"yandex-go-advanced/internal/logger"
	"yandex-go-advanced/internal/session"
	"yandex-go-advanced/internal/storage"
	pb "yandex-go-advanced/proto"

	"google.golang.org/grpc"
)

// logKeyError - error constant
const (
	logKeyError = "error"
)

func main() {
	cnf := config.Init()
	sgr := logger.Init()
	defer func() {
		err := sgr.Sync()
		if err != nil {
			log.Printf("failed to sync logger: %s", err.Error())
		}
	}()

	str, err := storage.Init(cnf)
	if err != nil {
		sgr.Errorw(
			"failed to init storage",
			logKeyError, err.Error(),
		)
		return
	}
	defer func() {
		err = str.Close()
		if err != nil {
			sgr.Errorw(
				"failed to close storage connection",
				logKeyError, err.Error(),
			)
		}
	}()

	ssp := &session.SessionProvider{
		Config: cnf,
	}

	listen, err := net.Listen("tcp", *cnf.Server)
	if err != nil {
		sgr.Errorw(
			"failed to listen",
			logKeyError, err.Error(),
		)
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterShortenerServiceServer(s, &grpc_handlers.HandlerProvider{
		Config:  cnf,
		Storage: str,
		Sugar:   sgr,
		Session: ssp,
	})

	sgr.Log(1, "server gRPC started in: ", *cnf.Server)

	if err := s.Serve(listen); err != nil {
		sgr.Errorw(
			"failed to serve",
			logKeyError, err.Error(),
		)
		log.Fatal(err)
	}
}
