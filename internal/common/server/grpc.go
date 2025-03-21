package server

import (
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"net"
)

// init 初始化 zap 全局日志器
func init() {
	config := zap.NewProductionConfig()
	zapLevel := viper.GetString("zap-level")
	switch zapLevel {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
	}
	logger, _ := config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
	zap.ReplaceGlobals(logger)
}

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
