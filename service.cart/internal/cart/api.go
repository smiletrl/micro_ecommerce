package cart

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	"net/http"
	"strconv"
)

// RegisterHandlers for handlers
func RegisterHandlers(r *echo.Group, service Service) {
	res := &resource{service}

	group := r.Group("/cart_item")

	group.GET("", res.Get)

	group.POST("", res.Create)

	group.DELETE("/:id", res.Delete)
}

type resource struct {
	service Service
}

type createRequest struct {
	Quantity  int   `db:"quantity"`
	ProductID int64 `db:"product_id"`
}

type createResponse struct {
	ID int64 `json:"id"`
}

func (r resource) Get(c echo.Context) error {
	return c.String(http.StatusOK, "succeed!")
}

func (r resource) Create(c echo.Context) error {
	req := new(createRequest)
	if err := c.Bind(req); err != nil {
		return errorsd.BadRequest(c, errors.Wrapf(errorsd.New("error creating request customer"), "error binding creating customer request: %s", err.Error()))
	}
	fmt.Printf("req is: %+v\n", req)
	customerID := int64(12)

	// get the product title, price, stock -- they could be retrieved from service product.
	// RPC call to service product.
	id, err := r.service.Create(c, customerID, req.ProductID, req.Quantity)
	return c.String(http.StatusOK, "succeed!")
	if err != nil {
		return errorsd.Abort(c, errors.Wrapf(errorsd.New("error creating customer"), "error creating customer: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, createResponse{
		ID: id,
	})
}

type deleteResponse struct {
	Data string `json:"data"`
}

func (r resource) Delete(c echo.Context) error {
	id := c.Param("id")
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errorsd.BadRequest(c, errors.Wrapf(errorsd.New("error getting request customer"), "error getting customer request: %s", err.Error()))
	}

	err = r.service.Delete(c, idInt64)
	if err != nil {
		return errorsd.Abort(c, errors.Wrapf(errorsd.New("error deleting customer"), "error deleting customer: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, deleteResponse{
		Data: "ok",
	})
}
