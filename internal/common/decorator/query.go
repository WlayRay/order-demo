package decorator

import (
	"context"
	"go.uber.org/zap"
)

// QueryHandler 定义一个接收查询的泛型类型Q,
// 然后返回一个泛型结果R
type QueryHandler[Q, R any] interface {
	Handle(context.Context, Q) (R, error)
}

func ApplyQueryDecorators[H, R any](handler QueryHandler[H, R], logger *zap.Logger, metricsClient MetricsClient) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		logger: logger,
		base: queryMetricsDecorator[H, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}
