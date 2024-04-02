package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nguyentrunghieu15/common-vcs-prj/apu/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config map[string]string

func ReadConfig(filename string) (Config, error) {
	// init with some bogus data
	config := Config{}
	if len(filename) == 0 {
		return config, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')

		// check if the line has = sign
		// and process the line. Ignore the rest.
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				// assign the config map
				config[key] = value
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}

func main() {
	properties_path := os.Getenv("PROPERTIES_PATH")
	config, err := ReadConfig(properties_path)
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	if config["enableswagger"] == "true" {
		e.Static("/static", "assets")
	}
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = auth.RegisterAuthServiceHandlerFromEndpoint(context.Background(), mux, "auth:3456", opts)

	if err != nil {
		log.Fatalln("Failed to mux server:", err)
	}
	e.Any("/*", echo.WrapHandler(mux)) // all HTTP requests starting with `/prefix` are handled by `grpc-gateway`

	e.Logger.Fatal(e.Start(":3000"))
}
