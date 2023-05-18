package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"os"
	"time"
)

type gRPC struct {
	addr   string
	server *grpc.Server
}

func (s *gRPC) handleShutdown(shutdown <-chan bool, done chan<- bool) {
	<-shutdown
	s.server.GracefulStop()
	log.Info().Msg("gRPC server stopped")
	done <- true
}

func (s *gRPC) start() {
	s.addr = ":50051"
	if os.Getenv("GRPC_PORT") != "" {
		s.addr = fmt.Sprintf(":%s", os.Getenv("GRPC_PORT"))
	}
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
	if err = s.server.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve gRPC")
	}
}

func (s *gRPC) setup() {
	opts := []grpc.ServerOption{
		grpc.Creds(insecure.NewCredentials()),
	}
	s.server = grpc.NewServer(opts...)
}

// NewServer initializes the gRPC service and the server object.
func NewServer(service *Server, shutdown <-chan bool, done chan<- bool) {
	s := &gRPC{}
	s.setup()
	service.RegisterGRPC(s.server)
	go s.start()
	go s.handleShutdown(shutdown, done)
	time.Sleep(1 * time.Second)
	log.Info().Msgf("Listening and serving gRPC on %s", s.addr)
}
