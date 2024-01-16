package grpcserver

import (
	"fmt"
	handlers "gRPCserver/internal/transport/handlers"
	"log/slog"
	"net"

	proto "github.com/Ivan010403/proto/protoc/go"
	"google.golang.org/grpc"
)

type Server struct {
	gRPCsrv *grpc.Server
	log     *slog.Logger
	port    int
}

func NewServer(logger *slog.Logger, port, wrkUploadGetFiles, wrkCheckFiles int, wrk handlers.FileWork) *Server {
	srv := grpc.NewServer()

	wrkUploadGet := make(chan struct{}, wrkUploadGetFiles)
	wrkCheck := make(chan struct{}, wrkCheckFiles)

	proto.RegisterCloudServer(srv, &handlers.CloudServer{
		ChanUploadGet: wrkUploadGet,
		ChanCheck:     wrkCheck,
		Worker:        wrk,
	})

	return &Server{
		gRPCsrv: srv,
		port:    port,
		log:     logger,
	}
}

func (s *Server) MustRun() {
	s.log.Info("starting gRPC server", slog.Int("port", s.port))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		panic(fmt.Sprintf("failed to listen port: %d", s.port))
	}

	if err := s.gRPCsrv.Serve(lis); err != nil {
		panic(fmt.Sprintf("failed to serve port: %d", s.port))
	}
}

func (s *Server) GracefulStop() {
	s.log.Info("stopping grpc server")

	s.gRPCsrv.GracefulStop()
}
