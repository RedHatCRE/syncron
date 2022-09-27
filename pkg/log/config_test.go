package log_test

import (
	"testing"

	"github.com/rhcre/syncron/pkg/log"
	"github.com/rhcre/syncron/test/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigureTestSuite struct {
	suite.Suite
}

func (suite *ConfigureTestSuite) TestLevelOnDebug() {
	parser := new(mocks.CLIParserMock)
	parser.On("GetDebug").Return(true)

	log.Configure(parser)

	assert.Equal(suite.T(), logrus.TraceLevel, logrus.GetLevel())
}

func (suite *ConfigureTestSuite) TestLevelOnNoDebug() {
	parser := new(mocks.CLIParserMock)
	parser.On("GetDebug").Return(false)

	log.Configure(parser)

	assert.Equal(suite.T(), logrus.InfoLevel, logrus.GetLevel())
}

func TestConfigureTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigureTestSuite))
}
