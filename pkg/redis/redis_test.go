package redis

import (
	"github.com/smiletrl/micro_ecommerce/pkg/config"
	"github.com/smiletrl/micro_ecommerce/pkg/constants"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTestDB(t *testing.T) {
	stage := os.Getenv(constants.Stage)
	if stage == "" {
		stage = constants.StageLocal
	}
	cfg, err := config.Load(stage)
	assert.NoError(t, err)
	client := Test(cfg)
	assert.NotNil(t, client)
}
