package context

import (
	"github.com/labstack/echo/v4"
)

// CustomContext is our customized context wrapper of echo context
// @todo remove this one. The default echo.Context is able to handle most cases
// already
type Context interface {
	GetCustomerID() int64
	SetCustomerID(id int64)
}

// Context is custom context
type context struct {
	echo.Context
	CustomerID int64
}

// New is to create new Context
func New(c echo.Context) Context {
	return context{Context: c}
}

func (c context) SetCustomerID(id int64) {
	c.CustomerID = id
}

func (c context) GetCustomerID() int64 {
	return c.CustomerID
}

func Mock(c echo.Context) Context {
	return context{
		Context:    c,
		CustomerID: int64(1),
	}
}
