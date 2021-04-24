package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/jwt"
	"net/http"
)

func CustomerMiddleware(jwtService jwt.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			customerID, err := jwtService.ParseCustomerToken(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			// set the customer id into context. We may want to set other info into context
			c.Set("customer_id", customerID)
			return next(c)
		}
	}
}

func AdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// @todo add admin jwt token
			return next(c)
		}
	}
}
