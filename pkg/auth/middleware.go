package auth

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/jwt"
	_ "net/http"
)

func CustomerMiddleware(jwtProvider jwt.Provider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//customerID, err := jwtProvider.ParseCustomerToken(c)
			//if err != nil {
			//return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			//}
			// set the customer id into context. We may want to set other info into context
			customerID := int64(12)
			c.Set(constants.AuthCustomerID, customerID)
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
