package decorator

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/context"
)

// MetricsClient is an interface for metrics clients.
type MetricsClient interface {
	Inc(key string, value int)
	Observe(key string, value float64)
}

type queryMetricsDecorator[C, R any] struct {
	base   QueryHandler[C, R]
	client MetricsClient
}

func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(cmd))
	defer func() {
		end := time.Since(start)
		q.client.Observe(fmt.Sprintf("%s.%s", actionName, "duration"), float64(end.Milliseconds()))
		if err == nil {
			q.client.Inc(fmt.Sprintf("%s.%s", actionName, "success"), 1)
		} else {
			q.client.Inc(fmt.Sprintf("%s.%s", actionName, "failure"), 1)
		}
	}()

	return q.base.Handle(ctx, cmd)
}

type commandMetricsDecorator[C, R any] struct {
	base   QueryHandler[C, R]
	client MetricsClient
}

func (q commandMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	start := time.Now()
	actionName := strings.ToLower(generateActionName(cmd))
	defer func() {
		end := time.Since(start)
		q.client.Observe(fmt.Sprintf("%s.%s", actionName, "duration"), float64(end.Milliseconds()))
		if err == nil {
			q.client.Inc(fmt.Sprintf("%s.%s", actionName, "success"), 1)
		} else {
			q.client.Inc(fmt.Sprintf("%s.%s", actionName, "failure"), 1)
		}
	}()

	return q.base.Handle(ctx, cmd)
}
