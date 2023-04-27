package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	app1v1 "github.com/Shitomo/play-with-chatgpt-4/pkg/connect/app1/v1"
	"github.com/Shitomo/play-with-chatgpt-4/pkg/connect/app1/v1/app1v1connect"
	"github.com/bufbuild/connect-go"
	otelconnect "github.com/bufbuild/connect-opentelemetry-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func NewExporter() (sdktrace.SpanExporter, error) {
	return otlptracegrpc.New(
		context.Background(),
		otlptracegrpc.WithInsecure(),
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

var tracer = otel.Tracer("app2/app2-service")

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
	url := os.Getenv("API_SERVER_URL")
	if url == "" {
		log.Fatal("API_SERVER_URL is not set")
	}

	client := app1v1connect.NewHelloServiceClient(http.DefaultClient, url, connect.WithInterceptors(otelconnect.NewInterceptor()))

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "app2/hello")
		defer span.End()
		log.Print("handle hello")

		res, err := client.Hello(ctx, connect.NewRequest(&app1v1.HelloRequest{}))
		if err != nil {
			log.Printf("failed to call hello: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write([]byte(fmt.Sprintf("%s, and Goodbye", res.Msg.Message)))
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
