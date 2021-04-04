package context

import (
	"github.com/labstack/echo"
)

// CustomContext is our customized context wrapper of echo context
type CustomContext interface {
	GetCustomerID() int64
	SetCustomerID(id int64)
}

// Context is custom context
type Context struct {
	echo.Context
	CustomerID int64
}

// New is to create new Context
func New(c echo.Context) *Context {
	return &Context{Context: c}
}

func (c *Context) SetCustomerID(id int64) {
	c.CustomerID = id
}

func (c *Context) GetUserID() int64 {
	return c.CustomerID
}

func NewMock(c echo.Context) *Context {
	return &Context{
		Context:    c,
		CustomerID: int64(1),
	}
}
