package server

import (
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

// RunGRPCServer 启动 gRPC 服务
func RunGRPCServer(serviceName string, registerServer func(server *grpc.Server)) {
	addr := viper.Sub(serviceName).GetString("grpc-addr")
	if addr == "" {
		zap.L().Warn("grpc-addr not found in config, using fallback-grpc-addr",
			zap.String("service", serviceName))
		addr = viper.GetString("fallback-grpc-addr")
		if addr == "" {
			zap.L().Fatal("both grpc-addr and fallback-grpc-addr are empty, cannot start server")
		}
	}
	RunGRPCServerOnAddr(addr, registerServer)
}

// RunGRPCServerOnAddr 在指定地址启动 gRPC 服务
func RunGRPCServerOnAddr(addr string, registerServer func(server *grpc.Server)) {
	grpcServer := grpc.NewServer(
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
