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
package log

import (
	"github.com/redhatcre/syncron/pkg/cli"
	"github.com/sirupsen/logrus"
)

// Configure sets up Logrus so that it behaves as the user describes
// on the CLI.
//
// At this moment, Configure supports the following flags from the CLI:
//   - '-d'
//   - '--debug'
func Configure(parser cli.CLIParser) error {
	// Define all actions needed to get Logrus ready to go
	setup := []func() error{
		// Set log level
		func() error {
			logrus.SetLevel(logrus.InfoLevel)

			if parser.GetDebug() {
				logrus.SetLevel(logrus.TraceLevel)
			}

			return nil
		},
	}

	// Run all steps defined above
	for _, step := range setup {
		err := step()

		// Halt if the setup failed
		if err != nil {
			return err
		}
	}

	return nil
}
