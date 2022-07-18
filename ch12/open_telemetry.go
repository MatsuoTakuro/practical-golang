package ch12

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("mux-server")

func initJaegerProvider(url string) (*sdktrace.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("otel demo"),
			attribute.String("environment", "local test"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

// curl localhost:8080/users/123
func openTelemetry() {
	tp1, err := initJaegerProvider("http://localhost:14268/api/traces")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp1.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down jaeger tracer provider: %v", err)
		}
	}()
	r := mux.NewRouter()
	r.Use(otelmux.Middleware("my-server"))

	r.HandleFunc("/users/{id:[0-9]+}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		name := getUser(r.Context(), id)
		time.Sleep(time.Microsecond)
		reply := fmt.Sprintf("user %s (id %s)\n", name, id)
		_, _ = w.Write(([]byte)(reply))
	}))
	http.Handle("/", r)
	log.Println("start listening at :8080")
	_ = http.ListenAndServe(":8080", nil)
}

func getUser(ctx context.Context, id string) string {
	time.Sleep(time.Microsecond)
	_, span := tracer.Start(ctx, "getUser", trace.WithAttributes(attribute.String("id", id)))
	defer span.End()
	time.Sleep(time.Microsecond)
	if id == "123" {
		return "otelmux tester"
	}
	return "unknown"
}
