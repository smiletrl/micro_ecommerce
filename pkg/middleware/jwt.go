package middleware

import (
	"github.com/labstack/echo"
	contextd "github.com/smiletrl/micro_ecommerce/pkg/context"
	"github.com/smiletrl/micro_ecommerce/pkg/jwt"
	"net/http"
	"reflect"
)

// JWTAuthMiddleware jwt middleware. Now all ports use this same middleware. Later we
// may want to use different middleware for different ports.
func JWTAuthMiddleware(skipRoutes []string, jwtService jwt.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//  skip the route
			if pathInArray(c.Path(), skipRoutes) {
				return next(c)
			}
			customerID, err := jwtService.ParseToken(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			// set the used id into context. We may want to set other info into context
			cc := c.(*contextd.Context)
			cc.SetCustomerID(customerID)
			return next(cc)
		}
	}
}

func pathInArray(val string, array interface{}) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				return true
			}
		}
	}
	return false
}
