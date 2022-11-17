// Copyright 2022 Red Hat, Inc.
// All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.
package log_test

import (
	"testing"

	"github.com/redhatcre/syncron/pkg/log"
	"github.com/redhatcre/syncron/test/mocks"
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

	err := log.Configure(parser)
	if err != nil {
		logrus.Error("An error occured configuring logs", err)
	}

	assert.Equal(suite.T(), logrus.TraceLevel, logrus.GetLevel())
}

// TestLevelOnNoDebug calls log.Configure, checking for logrus.InfoLevel to be
// Logrus's level when '-d' is not present on the CLI.
func (suite *ConfigureTestSuite) TestLevelOnNoDebug() {
	parser := new(mocks.CLIParserMock)
	parser.On("GetDebug").Return(false)

	err := log.Configure(parser)
	if err != nil {
		logrus.Error("An error occured configuring logs", err)
	}
	assert.Equal(suite.T(), logrus.InfoLevel, logrus.GetLevel())
}

// TestConfigureTestSuite takes care of running all tests
// on ConfigureTestSuite.
func TestConfigureTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigureTestSuite))
}
