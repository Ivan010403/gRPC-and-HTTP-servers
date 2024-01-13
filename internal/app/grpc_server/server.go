package grpcserver

import (
	"fmt"
	handlers "gRPCserver/internal/transport/grpc"
	"log"
	"net"

	proto "github.com/Ivan010403/proto/protoc/go"
	"google.golang.org/grpc"
)

type Server struct {
	gRPCsrv *grpc.Server
	port    int
}

func NewServer(port int) *Server {
	srv := grpc.NewServer()

	proto.RegisterCloudServer(srv, &handlers.StreamHandler{})
	return &Server{
		gRPCsrv: srv,
		port:    port,
	}
}

func (s *Server) MustRun() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 4545))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := s.gRPCsrv.Serve(lis); err != nil {
		panic("can not run server")
	}
}

//TODO: add graceful shutdown
