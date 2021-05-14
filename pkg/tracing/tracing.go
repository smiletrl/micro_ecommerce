package tracing

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

type Provider interface {
	// SetupTracer creates a new Jaeger tracer.
	SetupTracer(serviceName string, c config.Config) (io.Closer, error)

	// Middleware starts a root trace for each request.
	Middleware(log logger.Provider) echo.MiddlewareFunc

	StartSpan(c context.Context, operationName string) (opentracing.Span, context.Context)

	// finsh span, primarily for mock purpose
	FinishSpan(span opentracing.Span)
}

type provider struct{}

func NewProvider() Provider {
	return provider{}
}

func (p provider) SetupTracer(serviceName string, c config.Config) (io.Closer, error) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          false,
			CollectorEndpoint: c.TracingEndpoint,
		},
	}

	closer, err := cfg.InitGlobalTracer(
		fmt.Sprintf("%s.%s", serviceName, c.Stage),
		jaegercfg.Logger(jaeger.StdLogger),
	)

	return closer, err
}

func (p provider) Middleware(log logger.Provider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()

			path := req.URL.Path

			operationName := req.Method + " " + path
			var (
				span opentracing.Span
				ctx  context.Context
			)

			if req.URL.Path != "/health" {
				// create a root span for this request
				span, ctx = p.StartSpan(c.Request().Context(), operationName)
				defer p.FinishSpan(span)

				r := c.Request().WithContext(ctx)
				c.SetRequest(r)
			}

			if err = next(c); err != nil {
				c.Error(err)
			}

			if req.URL.Path != "/health" {
				p.setSpanTags(req, res, c.RealIP(), span)
			}
			return
		}
	}
}

func (p provider) StartSpan(ctx context.Context, operationName string) (opentracing.Span, context.Context) {
	return opentracing.StartSpanFromContext(ctx, operationName)
}

func (p provider) FinishSpan(span opentracing.Span) {
	span.Finish()
}

func (p provider) setSpanTags(req *http.Request, res *echo.Response, ip string, span opentracing.Span) {
	ctxKeys := map[string]interface{}{
		"http.method":      req.Method,
		"http.url":         req.URL.String(),
		"http.status_code": res.Status,
	}
	for key, val := range ctxKeys {
		span.SetTag(key, val)
	}
}
