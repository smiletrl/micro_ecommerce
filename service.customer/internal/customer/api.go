package customer

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	errorsd "github.com/smiletrl/micro_ecommerce/pkg/errors"
	"net/http"
)

// RegisterHandlers for customer
func RegisterHandlers(r *echo.Group, service Service) {
	res := &resource{service}

	customerGroup := r.Group("")

	customerGroup.GET("/customer/:id", res.Get)

	customerGroup.POST("/customer/:id", res.Create)

	customerGroup.PUT("/customer/:id", res.Update)

	customerGroup.DELETE("/customer/:id", res.Delete)
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
	req := new(getRequest)
	if err := c.Bind(req); err != nil {
		return errorsd.BadRequest(c, errors.Wrapf(errorsd.New("error getting request customer id"), "error binding getting customer request: %s", err.Error()))
	}
	cus, err := r.service.Get(c, req.ID)
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
	ID        int64  `json:"id"`
	Email     string `db:"email"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

type updateResponse struct {
	Data string `json:"data"`
}

func (r resource) Update(c echo.Context) error {
	req := new(updateRequest)
	if err := c.Bind(req); err != nil {
		return errorsd.BadRequest(c, errors.Wrapf(errorsd.New("error updating request customer id"), "error binding updating customer request: %s", err.Error()))
	}
	err := r.service.Update(c, req.ID, req.Email, req.FirstName, req.LastName)
	if err != nil {
		return errorsd.Abort(c, errors.Wrapf(errorsd.New("error updating customer"), "error updating customer: %s", err.Error()))
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
	req := new(deleteRequest)
	if err := c.Bind(req); err != nil {
		return errorsd.BadRequest(c, errors.Wrapf(errorsd.New("error deleting request customer id"), "error binding deleting customer request: %s", err.Error()))
	}
	err := r.service.Delete(c, req.ID)
	if err != nil {
		return errorsd.Abort(c, errors.Wrapf(errorsd.New("error deleting customer"), "error deleting customer: %s", err.Error()))
	}
	return c.JSON(http.StatusOK, deleteResponse{
		Data: "ok",
	})
}
