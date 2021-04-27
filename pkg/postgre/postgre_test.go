package postgre

import (
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitDB(t *testing.T) {
	stage := os.Getenv(constants.Stage)
	if stage == "" {
		stage = constants.StageLocal
	}
	cfg, err := config.Load(stage)
	assert.NoError(t, err)
	db, err := InitDB(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, db)
}
