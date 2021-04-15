package product

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"

	"net/http"
	"strconv"
)

// RegisterHandlers for product
func RegisterHandlers(r *echo.Group, service Service) {
	res := &resource{service}

	productGroup := r.Group("/product")

	productGroup.GET("/:id", res.Get)

	productGroup.POST("", res.Create)

	productGroup.PUT("/:id", res.Update)

	productGroup.DELETE("/:id", res.Delete)
}

type resource struct {
	service Service
}

type getRequest struct {
	ID int64 `json:"id"`
}

type getResponse struct {
	Data product `json:"data"`
}

func (r resource) Get(c echo.Context) error {
	id := c.Param("id")
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errorsd.BadRequest(c, errors.Wrapf(errorsd.New("error getting request customer"), "error getting customer request: %s", err.Error()))
	}
	cus, err := r.service.Get(c, idInt64)
	if err != nil {
		return errorsd.Abort(c, errors.Wrapf(errorsd.New("error getting product"), "error getting product: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, getResponse{
		Data: cus,
	})
}

type createRequest struct {
	Title  string `json:"title"`
	Amount int    `json:"amount"`
	Stock  int    `json:"stock"`
}

type createResponse struct {
	ID int64 `json:"id"`
}

func (r resource) Create(c echo.Context) error {
	req := new(createRequest)
	if err := c.Bind(req); err != nil {
		return errorsd.BadRequest(c, errors.Wrapf(errorsd.New("error creating request product"), "error binding creating product request: %s", err.Error()))
	}
	id, err := r.service.Create(c, req.Title, req.Amount, req.Stock)
	if err != nil {
		return errorsd.Abort(c, errors.Wrapf(errorsd.New("error creating product"), "error creating product: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, createResponse{
		ID: id,
	})
}

type updateRequest struct {
	Title  string `json:"title"`
	Amount int    `json:"amount"`
	Stock  int    `json:"stock"`
}

type updateResponse struct {
	Data string `json:"data"`
}

func (r resource) Update(c echo.Context) error {
	id := c.Param("id")
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errorsd.BadRequest(c, errors.Wrapf(errorsd.New("error getting request customer"), "error getting customer request: %s", err.Error()))
	}
	req := new(updateRequest)
	if err := c.Bind(req); err != nil {
		return errorsd.BadRequest(c, errors.Wrapf(errorsd.New("error updating request product id"), "error binding updating product request: %s", err.Error()))
	}
	err = r.service.Update(c, idInt64, req.Title, req.Amount, req.Stock)
	if err != nil {
		return errorsd.Abort(c, errors.Wrapf(errorsd.New("error updating product"), "error updating product: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, updateResponse{
		Data: "ok",
	})
}

type deleteRequest struct {
	ID int64 `json:"id"`
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
		return errorsd.Abort(c, errors.Wrapf(errorsd.New("error deleting product"), "error deleting product: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, deleteResponse{
		Data: "ok",
	})
}
