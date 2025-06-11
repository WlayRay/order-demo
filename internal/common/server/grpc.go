package server

import (
	"net"

	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// RunGRPCServer starts a gRPC server.
func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
	addr := viper.Sub(serviceName).GetString("grpc-addr")
	if addr == "" {
		zap.L().Fatal("grpc-addr not found in config", zap.String("service", serviceName))
	}
	runGRPCServerOnAddr(addr, registerServer)
}

// runGRPCServerOnAddr starts a gRPC server on a specified address.
func runGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			grpcTags.UnaryServerInterceptor(grpcTags.WithFieldExtractor(grpcTags.CodeGenRequestFieldExtractor)),
			grpcZap.UnaryServerInterceptor(zap.L(), grpcZap.WithMessageProducer(grpcZap.DefaultMessageProducer)),
		),
		grpc.ChainStreamInterceptor(
			grpcTags.StreamServerInterceptor(grpcTags.WithFieldExtractor(grpcTags.CodeGenRequestFieldExtractor)),
			grpcZap.StreamServerInterceptor(zap.L(), grpcZap.WithMessageProducer(grpcZap.DefaultMessageProducer)),
		),
	)
	registerServer(grpcServer)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		zap.L().Fatal("failed to listen", zap.String("addr", addr), zap.Error(err))
	}

	zap.L().Info("starting gRPC server", zap.String("addr", addr))
	if err := grpcServer.Serve(listen); err != nil {
		zap.L().Fatal("failed to serve", zap.Error(err))
	}
}
