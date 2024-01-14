package app

import (
	grpcserver "gRPCserver/internal/app/grpc_server"
	cloudStorage "gRPCserver/internal/services"
	"log/slog"
)

type App struct {
	GRPCsrv *grpcserver.Server
}

func NewApp(logger *slog.Logger, port, wrkSaveDelete, wrkCheckFiles int) *App {
	//TODO: postgres

	//TODO:

	cloud := cloudStorage.NewCloud(logger)

	srv := grpcserver.NewServer(logger, port, wrkSaveDelete, wrkCheckFiles, cloud)

	return &App{
		GRPCsrv: srv,
	}
}
