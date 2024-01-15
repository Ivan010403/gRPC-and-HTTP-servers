package app

import (
	grpcserver "gRPCserver/internal/app/grpc_server"
	"gRPCserver/internal/config"
	cloudStorage "gRPCserver/internal/services"
	"gRPCserver/internal/storage/postgres"
	"log/slog"
)

type App struct {
	GRPCsrv *grpcserver.Server
}

func NewApp(logger *slog.Logger, cfggrpc config.GRPC_server, cfg config.DataBase) *App {

	storage, err := postgres.New(cfg.Host, cfg.User, cfg.Password, cfg.Dbname, cfg.Port)
	if err != nil {
		logger.Error("creation db error", slog.Any("err", err))
		return nil
	}

	logger.Info("database has been created successfully")

	cloud := cloudStorage.NewCloud(logger, storage)

	srv := grpcserver.NewServer(logger, cfggrpc.Port, cfggrpc.MaxReadWriteConn, cfggrpc.MaxCheckConn, cloud)

	return &App{
		GRPCsrv: srv,
	}
}
