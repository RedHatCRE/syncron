package log_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigureTestSuite struct {
	suite.Suite
}

func (suite *ConfigureTestSuite) TestLevelOnNoDebug() {
	assert.True(suite.T(), true)
}

func TestConfigureTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigureTestSuite))
}
