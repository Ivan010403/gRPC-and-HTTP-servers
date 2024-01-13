package main

import "gRPCserver/internal/app"

func main() {
	application := app.NewApp(4545)

	application.GRPCsrv.MustRun()
}
