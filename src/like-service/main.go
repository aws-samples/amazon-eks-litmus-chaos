package main

import (
	"context"
	"fmt"
	"like-service/common"
	"like-service/handler"
	"like-service/repository"
	"like-service/utils"
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

var logger = log.New(os.Stdout, "like-service-api ", log.LstdFlags)

func bootstrapConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config/")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatalf("fatal error config file: %w", err)
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
	lockAddr := viper.GetString("lockAddr")
	dbHost := viper.GetString("dbHost")
	dbPort := viper.GetString("dbPort")
	dbUser := viper.GetString("dbUser")
	dbName := viper.GetString("dbName")
	dbPassword := viper.GetString("dbPassword")
	serviceName := viper.GetString("serviceName")
	serviceNamespace := viper.GetString("serviceNamespace")

	shutdown := initProviderProm(serviceName, serviceNamespace)
	defer shutdown()

	locker, err := utils.NewDatabase(lockAddr) // mock
	if err != nil {
		log.Fatalf("Failed to connect to redis lock: %s", err.Error())
	}
	log.Printf("Connected to redis...")

	database, err := repository.NewDatabase(dbHost, dbUser, dbPassword, dbName, dbPort) // mock
	if err != nil {
		log.Fatalf("Failed to connect to postgres: %s", err.Error())
	}
	log.Printf("Connected to postgres...")

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to get hostname: %s", err.Error())
	}
	log.Printf("Hostname: %s", hostname)
	common.Hostname = hostname

	// Create handlers
	ch := handler.NewLike(logger, database, locker)

	sm := http.NewServeMux()
	sm.Handle(fmt.Sprintf("%s", baseApiPath), ch)
	sm.HandleFunc(fmt.Sprintf("/healthz"), handler.GetHealthz)
	sm.Handle("/metrics", promhttp.Handler())
	corsHandler := cors.Default().Handler(sm)

	handler := otelhttp.NewHandler(corsHandler, "like")

	s := &http.Server{
		Addr:         fmt.Sprintf("%s", apiListenAddr),
		Handler:      handler,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	logger.Printf("Starting HTTP server...")
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()
	logger.Printf("Started HTTP server...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sig := <-c
	logger.Println("Received terminate, graceful shutdown", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		db, _ := database.Client.DB()
		db.Close()
		locker.Client.Close()
		cancel()
	}()

	if err := s.Shutdown(ctx); err != nil {
		logger.Fatalf("Server Shutdown Failed:%+v", err)
	}
}
