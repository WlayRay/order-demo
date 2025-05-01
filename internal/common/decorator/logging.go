package decorator

import (
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"strings"
)

type queryLoggingDecorator[C, R any] struct {
	logger *zap.Logger
	base   QueryHandler[C, R]
}

func (q queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	logger := q.logger.With(
		zap.Any("query", generateActionName(cmd)),
		zap.String("query_body", fmt.Sprintf("%+v", cmd)),
	)
	// logger.Debug("Executing query")
	defer func() {
		if err != nil {
			logger.Error("Query failed", zap.Error(err))
		} else {
			logger.Info("Query executed successfully")
		}
	}()

	return q.base.Handle(ctx, cmd)
}

type commandLoggingDecorator[C, R any] struct {
	logger *zap.Logger
	base   QueryHandler[C, R]
}

func (q commandLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	logger := q.logger.With(
		zap.Any("command", generateActionName(cmd)),
		zap.String("command_body", fmt.Sprintf("%+v", cmd)),
	)
	// logger.Debug("Executing command")
	defer func() {
		if err != nil {
			logger.Error("Command failed", zap.Error(err))
		} else {
			logger.Info("Command executed successfully")
		}
	}()

	return q.base.Handle(ctx, cmd)
}


func generateActionName(cmd any) string {
	return strings.Split(fmt.Sprintf("%T", cmd), ".")[1]
}
