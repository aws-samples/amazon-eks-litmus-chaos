package main

import (
	"context"
	"counter-service/common"
	"counter-service/handler"
	"counter-service/repository"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	_ "go.uber.org/automaxprocs"
)

var l = log.New(os.Stdout, "counter-service-api ", log.LstdFlags)

func bootstrapConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config/")
	err := viper.ReadInConfig()
	if err != nil {
		l.Fatalf("fatal error config file: %w", err)
	}
}

func initProviderProm(serviceName string, serviceNamespace string) func() {
	ctx := context.Background()

	exporter, err := prometheus.New()
	if err != nil {
		panic(err)
	}

	resources, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceNamespaceKey.String(serviceNamespace),
		),
	)

	provider := metric.NewMeterProvider(
		metric.WithResource(resources),
		metric.WithReader(exporter))

	global.SetMeterProvider(provider)

	return func() {
		cxt, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		if err := provider.Shutdown(cxt); err != nil {
			otel.Handle(err)
		}
	}
}

func main() {
	// Read config
	bootstrapConfig()
	baseApiPath := viper.GetString("baseApiPath")
	apiListenAddr := viper.GetString("apiListenAddr")
	redisAddr := viper.GetString("redisAddr")
	serviceName := viper.GetString("serviceName")
	serviceNamespace := viper.GetString("serviceNamespace")

	shutdown := initProviderProm(serviceName, serviceNamespace)
	defer shutdown()

	database, err := repository.NewDatabase(redisAddr) // mock
	if err != nil {
		log.Fatalf("Failed to connect to redis: %s", err.Error())
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to get hostname: %s", err.Error())
	}
	log.Printf("Hostname: %s", hostname)
	common.Hostname = hostname

	// Create handlers
	ch := handler.NewCounter(l, database)

	sm := http.NewServeMux()
	sm.Handle(fmt.Sprintf("%s/count", baseApiPath), ch)
	sm.HandleFunc(fmt.Sprintf("/healthz"), handler.GetHealthz)
	sm.Handle("/metrics", promhttp.Handler())
	corsHandler := cors.Default().Handler(sm)

	handler := otelhttp.NewHandler(corsHandler, "counter")

	s := &http.Server{
		Addr:         fmt.Sprintf("%s", apiListenAddr),
		Handler:      handler,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	l.Printf("Starting HTTP server...")
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()
	l.Printf("Started HTTP server...")
	log.Printf("serving metrics at %s/metrics", apiListenAddr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	l.Println("Received terminate, graceful shutdown", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		database.Client.Close()
		cancel()
	}()

	if err := s.Shutdown(ctx); err != nil {
		l.Fatalf("Server Shutdown Failed:%+v", err)
	}
}
