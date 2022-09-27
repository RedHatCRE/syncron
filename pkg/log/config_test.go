package log_test

import (
	"testing"

	"github.com/rhcre/syncron/pkg/log"
	"github.com/rhcre/syncron/test/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ConfigureTestSuite groups all tests that target Configure.
type ConfigureTestSuite struct {
	suite.Suite
}

// TestLevelOnDebug calls log.Configure, checking for logrus.TraceLevel to be
// Logrus's level when '-d' is present on the CLI.
func (suite *ConfigureTestSuite) TestLevelOnDebug() {
	parser := new(mocks.CLIParserMock)
	parser.On("GetDebug").Return(true)

	log.Configure(parser)

	assert.Equal(suite.T(), logrus.TraceLevel, logrus.GetLevel())
}

// TestLevelOnNoDebug calls log.Configure, checking for logrus.InfoLevel to be
// Logrus's level when '-d' is not present on the CLI.
func (suite *ConfigureTestSuite) TestLevelOnNoDebug() {
	parser := new(mocks.CLIParserMock)
	parser.On("GetDebug").Return(false)

	log.Configure(parser)

	assert.Equal(suite.T(), logrus.InfoLevel, logrus.GetLevel())
}

// TestConfigureTestSuite takes care of running all tests
// on ConfigureTestSuite.
func TestConfigureTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigureTestSuite))
}
