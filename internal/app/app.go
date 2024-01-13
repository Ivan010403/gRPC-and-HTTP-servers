package app

import (
	grpcserver "gRPCserver/internal/app/grpc_server"
)

type App struct {
	GRPCsrv *grpcserver.Server
	Port    int
}

func NewApp(port int) *App {
	//TODO: postgres

	//TODO:

	srv := grpcserver.NewServer(port)

	return &App{
		GRPCsrv: srv,
	}
}
