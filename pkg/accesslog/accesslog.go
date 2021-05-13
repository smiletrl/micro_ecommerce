package accesslog

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"strconv"
	"time"
)

func Middleware(logger logger.Provider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			res := c.Response()
			start := time.Now()
			if err = next(c); err != nil {
				c.Error(err)
			}
			if res.Status != 200 || req.URL.Path != "/health" {
				var sh, _ = time.LoadLocation("Asia/Shanghai")
				var layout = "2006-01-02 15:04:05"
				stop := time.Now()
				cl := req.Header.Get(echo.HeaderContentLength)
				if cl == "" {
					cl = "0"
				}
				logger.Infow("http request",
					"bytes_in", cl,
					"bytes_out", strconv.FormatInt(res.Size, 10),
					"user_agent", req.UserAgent(),
					"remote_ip", c.RealIP(),
					"http.host", req.Host,
					"http.status", strconv.Itoa(res.Status),
					"http.uri", req.RequestURI,
					"http.method", req.Method,
					"latency", stop.Sub(start).String(),
					"time", start.In(sh).Format(layout),
				)
			}
			return
		}
	}
}
