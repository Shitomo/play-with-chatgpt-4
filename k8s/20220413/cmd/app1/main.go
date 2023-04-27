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
	otelconnect "github.com/bufbuild/connect-opentelemetry-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func NewExporter() (sdktrace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
		stdouttrace.WithWriter(os.Stderr),
	)
}

func NewResource(name, version string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(name),
		semconv.ServiceVersionKey.String(version),
	)
}

func SetupTraceProvider(shutdownTimeout time.Duration) (func(), error) {
	exporter, err := NewExporter()
	if err != nil {
		return nil, err
	}

	reource := NewResource("example-service", "1.0.0")
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(reource),
	)
	otel.SetTracerProvider(tracerProvider)

	cleanup := func() {
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown tracer provider: %v", err)
		}
	}
	return cleanup, nil
}

var tracer = otel.Tracer("app1/app1-service")

var _ app1v1connect.HelloServiceHandler = (*App1Server)(nil)

type App1Server struct{}

func (a *App1Server) Hello(
	ctx context.Context,
	req *connect.Request[app1v1.HelloRequest],
) (*connect.Response[app1v1.HelloResponse], error) {
	_, span := tracer.Start(ctx, "app2/hello")
	defer span.End()
	log.Println("handle hello")

	return connect.NewResponse(&app1v1.HelloResponse{
		Message: "Hello World. Hello Service",
	}), nil
}

func main() {
	cleanup, err := SetupTraceProvider(10 * time.Second)
	if err != nil {
		log.Fatal(err)
	}

	defer cleanup()

	// 後続のサービスにつなげるためにpropagaterを追加
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	app1Server := &App1Server{}

	mux := http.NewServeMux()

	mux.Handle(app1v1connect.NewHelloServiceHandler(app1Server, connect.WithInterceptors(otelconnect.NewInterceptor())))

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
