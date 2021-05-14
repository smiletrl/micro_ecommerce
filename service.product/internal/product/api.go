package product

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/auth"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"

	"net/http"
)

// RegisterHandlers for product
func RegisterHandlers(r *echo.Group, service Service) {
	res := &resource{service}

	adminGroup := r.Group("/product")
	adminGroup.Use(auth.AdminMiddleware())

	adminGroup.GET("/:id", res.Get)

	adminGroup.POST("", res.Create)

	adminGroup.PUT("/:id", res.Update)

	adminGroup.DELETE("/:id", res.Delete)
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
	ctx := c.Request().Context()
	cus, err := r.service.Get(ctx, c.Param("id"))
	if err != nil {
		return errorsd.Abort(c, errors.Wrapf(errorsd.New("error getting product"), "error getting product: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, getResponse{
		Data: cus,
	})
}

type createRequest struct {
	Title    string        `json:"title"`
	Body     string        `json:"body"`
	Category string        `json:"category"`
	Assets   productAssets `json:"assets"`
	Variants variantConfig `json:"variants"`
}

type createResponse struct {
	Data ID `json:"data"`
}

func (r resource) Create(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(createRequest)
	if err := c.Bind(req); err != nil {
		return errorsd.BadRequest(c, err)
	}

	id, err := r.service.Create(ctx, *req)
	if err != nil {
		return errorsd.Abort(c, errors.Wrapf(errorsd.New("error creating product"), "error creating product: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, createResponse{
		Data: ID{ID: id},
	})
}

type updateRequest struct {
	Title    string        `json:"title"`
	Body     string        `json:"body"`
	Category string        `json:"category"`
	Assets   productAssets `json:"assets"`
	Variants variantConfig `json:"variants"`
}

type updateResponse struct {
	Data ID `json:"data"`
}

func (r resource) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	req := new(updateRequest)
	if err := c.Bind(req); err != nil {
		return errorsd.BadRequest(c, errors.Wrapf(errorsd.New("error updating request product id"), "error binding updating product request: %s", err.Error()))
	}
	err := r.service.Update(ctx, id, *req)
	if err != nil {
		return errorsd.Abort(c, errors.Wrapf(errorsd.New("error updating product"), "error updating product: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, updateResponse{
		Data: ID{ID: id},
	})
}

type deleteRequest struct {
	ID int64 `json:"id"`
}

type deleteResponse struct {
	Data string `json:"data"`
}

func (r resource) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	err := r.service.Delete(ctx, c.Param("id"))
	if err != nil {
		return errorsd.Abort(c, errors.Wrapf(errorsd.New("error deleting product"), "error deleting product: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, deleteResponse{
		Data: "ok",
	})
}
