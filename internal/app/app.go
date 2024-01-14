package app

import (
	grpcserver "gRPCserver/internal/app/grpc_server"
	"log/slog"
)

type App struct {
	GRPCsrv *grpcserver.Server
}

func NewApp(logger *slog.Logger, port, wrkSaveDelete, wrkCheckFiles int) *App {
	//TODO: postgres

	//TODO:

	srv := grpcserver.NewServer(logger, port, wrkSaveDelete, wrkCheckFiles)

	return &App{
		GRPCsrv: srv,
	}
}
