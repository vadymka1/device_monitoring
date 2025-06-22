package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"device-monitoring/internal/logger"
	proto "device-monitoring/internal/proto"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	grpcHandler "device-monitoring/internal/grpc"
	httpHandler2 "device-monitoring/internal/handlers/http"

	"device-monitoring/internal/repositories"
	"device-monitoring/internal/services"
)

func main() {
	log := logger.New()
	log.Info("Starting application...")

	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=monitor sslmode=disable")
	if err != nil {
		log.Error("failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := repositories.NewDeviceRepository(db, log)
	monitor := services.NewMonitorService(repo, log)
	service := services.NewDeviceService(*repo, monitor, log)
	httpHandler := httpHandler2.NewDeviceHandler(service, monitor, log)

	r := chi.NewRouter()
	httpHandler2.SetupRoutes(r, httpHandler)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	grpcServer := grpc.NewServer()
	grpcSrv := grpcHandler.NewMonitoringServer(repo, monitor, log)
	proto.RegisterMonitoringServer(grpcServer, grpcSrv)

	go func() {
		log.Info("🚀 Starting HTTP server on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server failed: %v", err)
			os.Exit(1)
		}
	}()

	go func() {
		listener, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Error("Failed to listen on port 50051: %v", err)
			os.Exit(1)
		}
		log.Info("🚀 Starting gRPC server on :50051")
		if err := grpcServer.Serve(listener); err != nil {
			log.Error("gRPC server failed: %v", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Gracefully shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("HTTP shutdown failed: %v", err)
		os.Exit(1)
	}

	grpcServer.GracefulStop()
	log.Info("Servers stopped.")
}
