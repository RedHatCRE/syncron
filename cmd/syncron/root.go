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

package cmd

import (
	"github.com/redhatcre/syncron/pkg/cli"
	"github.com/redhatcre/syncron/pkg/log"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Version:           "1.0.0",
	Use:               "syncron",
	Short:             "Syncron - Easily download files from s3 buckets",
	Example:           "syncron download sosreports --days 10",
	PersistentPreRunE: onPersistentPreRun,
}

func Execute() error {
	return root.Execute()
}

func init() {
	root.PersistentFlags().BoolP(
		cli.Debug,
		cli.D,
		false,
		"Turn on debug mode",
	)
}

func onPersistentPreRun(cmd *cobra.Command, args []string) error {
	// Define actions taken to get Syncron ready before commands are executed
	setup := []func() error{
		// Set up logging
		func() error {
			parser := cli.NewParserForCobra(cmd, args)
			return log.Configure(parser)
		},
	}
	for _, step := range setup {
		if err := step(); err != nil {
			return err
		}
	}
	return nil
}
