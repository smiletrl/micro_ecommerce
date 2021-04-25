package customer

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	"net/http"
	"strconv"
)

// RegisterHandlers for customer
func RegisterHandlers(r *echo.Group, service Service) {
	res := &resource{service}

	customerGroup := r.Group("/customer")

	customerGroup.GET("/:id", res.Get)

	customerGroup.POST("", res.Create)

	customerGroup.PUT("/:id", res.Update)

	customerGroup.DELETE("/:id", res.Delete)
}

type resource struct {
	service Service
}

type getRequest struct {
	ID int64 `json:"id"`
}

type getResponse struct {
	Data customer `json:"data"`
}

func (r resource) Get(c echo.Context) error {
	id := c.Param("id")
	idInt64, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return errorsd.BadRequest(c, errors.Wrapf(errorsd.New("error getting request customer"), "error getting customer request: %s", err.Error()))
	}
	cus, err := r.service.Get(c, idInt64)
	if err != nil {
		return errorsd.Abort(c, errors.Wrapf(errorsd.New("error getting customer"), "error getting customer: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, getResponse{
		Data: cus,
	})
}

type createRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type createResponse struct {
	ID int64 `json:"id"`
}

func (r resource) Create(c echo.Context) error {
	req := new(createRequest)
	if err := c.Bind(req); err != nil {
		return errorsd.BadRequest(c, errors.Wrapf(errorsd.New("error creating request customer"), "error binding creating customer request: %s", err.Error()))
	}
	id, err := r.service.Create(c, req.Email, req.FirstName, req.LastName)
	if err != nil {
		return errorsd.Abort(c, errors.Wrapf(errorsd.New("error creating customer"), "error creating customer: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, createResponse{
		ID: id,
	})
}

type updateRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
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
		return errorsd.BadRequest(c, errors.Wrapf(errorsd.New("error updating request customer id"), "error binding updating customer request: %s", err.Error()))
	}
	err = r.service.Update(c, idInt64, req.Email, req.FirstName, req.LastName)
	if err != nil {
		return errorsd.Abort(c, errors.Wrapf(errorsd.New("error updating customer"), "error updating customer: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, updateResponse{
		Data: "ok",
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
