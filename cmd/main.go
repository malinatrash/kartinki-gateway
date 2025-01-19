package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/malinatrash/kartinki-gateway/internal/config"
	"github.com/malinatrash/kartinki-gateway/internal/middleware"
	auth_service_pb "github.com/malinatrash/kartinki-proto/gen/go/auth_service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

func dialAuthService(ctx context.Context, host string, logger *slog.Logger) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	conn, err := grpc.Dial(host, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial auth service: %w", err)
	}

	timeout := time.After(10 * time.Second)
	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("auth service connection timeout")
		default:
			state := conn.GetState()
			logger.Info("auth service connection state", "state", state)

			if state == connectivity.Ready {
				return conn, nil
			}
			if !conn.WaitForStateChange(ctx, state) {
				return nil, fmt.Errorf("auth service is not available")
			}
		}
	}
}

func main() {
	logger := slog.Default()
	cfg := config.Load()

	// Log initial configuration
	logger.Info("starting gateway service with configuration",
		"gateway_host", cfg.GatewayHost,
		"gateway_port", cfg.GatewayPort,
		"auth_service", fmt.Sprintf("%s:%s", cfg.AuthHost, cfg.AuthPort))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize gRPC gateway multiplexer
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	// Setup HTTP server first
	handler := middleware.LoggerMiddleware(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			logger.Info("health check request received")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
			return
		}
		mux.ServeHTTP(w, r)
	}))

	addr := fmt.Sprintf("%s:%s", cfg.GatewayHost, cfg.GatewayPort)

	// Start server in goroutine
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		logger.Info("gateway server is starting", "address", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err, "address", addr)
			panic(err)
		}
	}()

	authHost := fmt.Sprintf("%s:%s", cfg.AuthHost, cfg.AuthPort)
	logger.Info("attempting to connect to auth service", "address", authHost)

	authConn, err := dialAuthService(ctx, authHost, logger)
	if err != nil {
		logger.Error("auth service connection failed",
			"error", err,
			"auth_host", cfg.AuthHost,
			"auth_port", cfg.AuthPort)
	} else {
		defer authConn.Close()

		logger.Info("successfully connected to auth service",
			"host", authHost,
			"state", authConn.GetState())

		// Register auth service handler
		opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
		logger.Info("registering auth service handler", "endpoint", authHost)

		if err := auth_service_pb.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, authHost, opts); err != nil {
			logger.Error("failed to register auth service handler",
				"error", err,
				"endpoint", authHost)
		} else {
			logger.Info("successfully registered auth service handler", "endpoint", authHost)
		}
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logger.Info("shutting down server...")

	// Shutdown server gracefully
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
	}

	logger.Info("server exited")
}
