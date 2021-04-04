package balance

import (
	"github.com/labstack/echo/v4"
	"github.com/smiletrl/micro_ecommerce/pkg/errors"
	"net/http"
)

// RegisterHandlers for balance
func RegisterHandlers(r *echo.Group, service Service) {
	res := resource{service}

	balanceGroup := r.Group("/balance")

	// This is a test purpose endpoint to add balance for customer
	balanceGroup.POST("/add", res.Add)
}

type resource struct {
	service Service
}

type addRequest struct {
	CustomerID int64 `json:"customer_id"`
	Balance    int   `json:"balance"`
}

type addResponse struct {
	Data string `json:"data"`
}

func (r resource) Add(c echo.Context) error {
	req := new(addRequest)
	if err := c.Bind(req); err != nil {
		return errors.BadRequest(c, err)
	}
	err := r.service.Add(c, req.CustomerID, req.Balance)
	if err != nil {
		return errors.Abort(c, err)
	}
	return c.JSON(http.StatusOK, addResponse{
		Data: "ok",
	})
}
