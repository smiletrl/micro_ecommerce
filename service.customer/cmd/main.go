package cmd

import (
	"github.com/labstack/echo"
	"github.com/smiletrl/micro_ecommerce/pkg/dbcontext"
	"github.com/smiletrl/micro_ecommerce/service.customer/internal/balance"
	"github.com/smiletrl/micro_ecommerce/service.customer/internal/customer"
)

// RegisterHandlersWechatMiniprogram registers staff handlers
func RegisterHandlersWechatMiniprogram(r *echo.Group, db dbcontext.DB, jwt jwt.Service, config *config.Config) {
	customer.RegisterHandlersWechatMiniProgram(r,
		customer.NewService(customer.NewRepository(db), jwt, config),
	)
}

func main() {
	customerRepo = customer.NewRepository()
	customerService = customer.NewService()
	customer.RegisterHandlers()
}
