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

//TODO: change name FileWork

func NewServer(logger *slog.Logger, port, wrkSaveDelete, wrkCheckFiles int, wrk handlers.FileWork) *Server {
	srv := grpc.NewServer()

	wrkSave := make(chan struct{}, wrkSaveDelete)
	wrkDelete := make(chan struct{}, wrkSaveDelete)
	wrkCheck := make(chan struct{}, wrkCheckFiles)

	proto.RegisterCloudServer(srv, &handlers.CloudServer{
		ChanSave:   wrkSave,
		ChanDelete: wrkDelete,
		ChanCheck:  wrkCheck,
		Worker:     wrk,
	})

	return &Server{
		gRPCsrv: srv,
		port:    port,
		log:     logger,
	}
}

func (s *Server) MustRun() {
	s.log.Info("starting gRPC server", slog.String("path", "app.grpc_server.server.MustRun"), slog.Int("port", s.port))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		panic(fmt.Sprintf("failed to listen port: %d", s.port))
	}

	if err := s.gRPCsrv.Serve(lis); err != nil {
		panic(fmt.Sprintf("failed to serve port: %d", s.port))
	}
}

func (s *Server) GracefulStop() {
	s.log.Info("stopping server", slog.String("path", "app.grpc_server.server.GracefulStop"))

	s.gRPCsrv.GracefulStop()
}
