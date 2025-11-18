package server

import (
	"context"
	"fmt"
	"net"

	"github.com/fatkulllin/gophkeeper/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func (server *Server) StartGRPC(ctx context.Context) error {
	listen, err := net.Listen("tcp", server.config.GRPCAddress)
	if err != nil {
		return fmt.Errorf("failed to listen on gRPC address: %w", err)
	}

	serverGRPC := grpc.NewServer()

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(serverGRPC, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	logger.Log.Info("gRPC server is running on", zap.String("address", server.config.GRPCAddress))

	go func() {
		<-ctx.Done()
		logger.Log.Info("Shutting down gRPC server gracefully")
		serverGRPC.GracefulStop()
	}()

	if err := serverGRPC.Serve(listen); err != nil {
		if err == grpc.ErrServerStopped {
			logger.Log.Info("gRPC server stopped gracefully")
			return nil
		}
		return fmt.Errorf("failed to serve gRPC: %w", err)
	}
	return nil
}
