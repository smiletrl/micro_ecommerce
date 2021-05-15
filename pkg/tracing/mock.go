package tracing

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"io"
)

type mockProvider struct{}

func NewMockProvider() (Provider, error) {
	return mockProvider{}, nil
}

func (m mockProvider) SetupTracer(serviceName string, c config.Config) (io.Closer, error) {
	return nil, nil
}

func (m mockProvider) Middleware(log logger.Provider) echo.MiddlewareFunc {
	return nil
}

func (m mockProvider) StartSpan(c context.Context, operationName string) (opentracing.Span, context.Context) {
	return nil, context.Background()
}

func (p mockProvider) FinishSpan(span opentracing.Span) {
}

func (p mockProvider) Close() {
}
