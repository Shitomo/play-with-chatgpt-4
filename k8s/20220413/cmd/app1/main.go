package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	app1v1 "github.com/Shitomo/play-with-chatgpt-4/pkg/connect/app1/v1"
	"github.com/Shitomo/play-with-chatgpt-4/pkg/connect/app1/v1/app1v1connect"
	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var _ app1v1connect.HelloServiceHandler = (*App1Server)(nil)

type App1Server struct{}

func (a *App1Server) Hello(
	ctx context.Context,
	req *connect.Request[app1v1.HelloRequest],
) (*connect.Response[app1v1.HelloResponse], error) {
	log.Println("handle hello")

	return connect.NewResponse(&app1v1.HelloResponse{
		Message: "Hello World",
	}), nil
}

func main() {
	app1Server := &App1Server{}

	mux := http.NewServeMux()

	mux.Handle(app1v1connect.NewHelloServiceHandler(app1Server))

	srv := &http.Server{
		Addr: fmt.Sprintf(":%v", 8080),
		Handler: h2c.NewHandler(
			mux,
			&http2.Server{},
		),
	}

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed with error: %s", err.Error())

		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)

	log.Printf("SIGNAL %d received, then shutting down...\n", <-quit)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("failed to gracefully shutdown: %v", err)
	}

	log.Printf("server shutdown")
}
