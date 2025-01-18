package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/malinatrash/kartinki-gateway/internal/config"
	pb "github.com/malinatrash/kartinki-proto/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger := slog.Default()
	cfg := config.Load()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterAuthServiceHandlerFromEndpoint(
		ctx,
		mux,
		fmt.Sprintf("%s:%s", cfg.AuthHost, cfg.AuthPort),
		opts,
	)
	if err != nil {
		panic(err)
	}

	logger.Info("server started on %s:%s", cfg.GatewayHost, cfg.GatewayPort)
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", cfg.GatewayHost, cfg.GatewayPort), mux)
}
