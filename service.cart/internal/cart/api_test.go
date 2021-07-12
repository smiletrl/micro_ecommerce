package cart

import (
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/smiletrl/micro_ecommerce/pkg/errors"
	"github.com/smiletrl/micro_ecommerce/pkg/jwt"
	"github.com/smiletrl/micro_ecommerce/pkg/logger"
	"github.com/smiletrl/micro_ecommerce/pkg/redis"
	"github.com/smiletrl/micro_ecommerce/pkg/test"
	"github.com/smiletrl/micro_ecommerce/pkg/tracing"
)

var tests = []test.APITestCase{
	// *** ROUTE: /cart
	// ** POST
	{"create cart correct", "POST", "/api/v1/cart", nil, `{"quantity": 12, "sku_id": "xxx"}`, http.StatusOK, `{"data":"ok"}`},
	{"create cart bad request", "POST", "/api/v1/cart", nil, `{"sku_id": "xxx"}`, http.StatusBadRequest, `{"code": "error", "message":"missing parameter quantity or sku"}`},
}

func TestAPI(t *testing.T) {
	e := echo.New()

	logger := logger.NewMockProvider()
	tracing, err := tracing.NewMockProvider()
	assert.NoError(t, err)

	// middleware
	e.Use(errors.Middleware(logger))
	group := e.Group("/api/v1")

	// config
	stage := os.Getenv(constants.Stage)
	if stage == "" {
		stage = constants.StageLocal
	}
	cfg, err := config.Load(stage)
	assert.NoError(t, err)

	// repository
	rdb := redis.NewMockProvider(cfg, tracing)
	repo := NewRepository(rdb, tracing)

	// service
	product := newMockProduct()
	jwt := jwt.NewMockProvider()
	mockService := NewService(repo, product, logger)

	RegisterHandlers(group, mockService, jwt)

	for _, tc := range tests {
		test.Endpoint(t, e, tc)
	}
}
