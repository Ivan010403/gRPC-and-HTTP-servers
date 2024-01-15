package main_test

import (
	"context"
	"gRPCserver/internal/config"
	"testing"
	"time"

	proto "github.com/Ivan010403/proto/protoc/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Test struct {
	*testing.T
	Cfg     *config.Config
	CloudCl proto.CloudClient
}

func New(t *testing.T) (context.Context, *Test) {
	t.Helper()

	cfg := config.ReadConfigFromPath("../config/local.yaml")

	ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second*20)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	conn, err := grpc.DialContext(context.Background(), "localhost:4545", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("connection failed")
	}
	return ctx, &Test{T: t, Cfg: cfg, CloudCl: proto.NewCloudClient(conn)}
}
