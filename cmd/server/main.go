package main

import (
	"gRPCserver/internal/app"
	"gRPCserver/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := createLogger()
	if logger == nil {
		panic("failed to create logger")
	}

	logger.Info("logger has been created successfully")

	cfg, err := config.ReadConfig()
	if err != nil {
		logger.Error("failed to read config", slog.Any("err", err))
		return
	}

	logger.Info("config has been read successfully")

	application := app.NewApp(logger, cfg.GRPC_server, cfg.DataBase)

	go application.GRPCsrv.MustRun()

	//graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan

	logger.Info("stopping the application")
	application.GRPCsrv.GracefulStop()
	logger.Info("application successfully stoped")

}

// TODO: add switch with env and different levels of logging
func createLogger() *slog.Logger {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return logger
}
