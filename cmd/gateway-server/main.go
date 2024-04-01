package main

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nguyentrunghieu15/common-vcs-prj/apu/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "assets")

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := auth.RegisterAuthServiceHandlerFromEndpoint(context.Background(), mux, "localhost:3456", opts)

	if err != nil {
		log.Fatalln("Failed to mux server:", err)
	}
	e.Any("/*", echo.WrapHandler(mux)) // all HTTP requests starting with `/prefix` are handled by `grpc-gateway`

	e.Logger.Fatal(e.Start(":3000"))
}
