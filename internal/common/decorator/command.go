package decorator

import (
	"context"
	"go.uber.org/zap"
)

type CommandHandler[C, R any] interface {
	Handle(context.Context, C) (R, error)
}

func ApplyCommandDecorators[C, R any](handler QueryHandler[C, R], logger *zap.Logger, metricsClient MetricsClient) QueryHandler[C, R] {
	return queryLoggingDecorator[C, R]{
		logger: logger,
		base: queryMetricsDecorator[C, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}
