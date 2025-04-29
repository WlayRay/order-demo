package decorator

import (
	"context"
	"go.uber.org/zap"
)

type CommandHandler[C, R any] interface {
	Handle(context.Context, C) (R, error)
}

// ApplyCommandDecorators applies decorators to a command handler.
func ApplyCommandDecorators[C, R any](handler QueryHandler[C, R], logger *zap.Logger, metricsClient MetricsClient) QueryHandler[C, R] {
	return commandLoggingDecorator[C, R]{
		logger: logger,
		base: commandMetricsDecorator[C, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}
