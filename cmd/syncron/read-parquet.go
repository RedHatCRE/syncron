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

// This is a reimplementation of the amazing work done by Mathew Topol.
// Thanks to his work, reading parquet files was made incredibly easy.

package cmd

import (
	"github.com/redhatcre/syncron/pkg/cli"
	"github.com/redhatcre/syncron/pkg/parquet_reader"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var parquet = &cobra.Command{
	Use:       "read-parquet",
	Short:     "Process parquet files",
	Args:      cobra.MatchAll(cobra.OnlyValidArgs),
	ValidArgs: []string{cli.Output, cli.File},
	RunE:      onRunParquet,
}

func init() {
	parquet.Flags().String(
		cli.Output,
		"-",
		"Specify output file for data.",
	)
	parquet.Flags().String(
		cli.File,
		"",
		"Specify file to read.",
	)

	// To be implemented in future PR.

	//parquet.Flags().Bool(
	//	"json",
	//	false,
	//	"Format output as JSON instead of text",
	//)
	//parquet.Flags().Bool(
	//	"csv",
	//	false,
	//	"Format output as CSV instead of text",
	//)
	//parquet.Flags().String(
	//	"columns",
	//	"",
	//	"Specify a subset of columns to print, comma delimited indexes.",
	//)
	root.AddCommand(parquet)
}

func onRunParquet(cmd *cobra.Command, args []string) error {
	output, _ := cmd.Flags().GetString(cli.Output)
	file, _ := cmd.Flags().GetString(cli.File)
	logrus.Info("File exported to: ", output)
	parquet_reader.ReadParquet(output, file)
	return nil
}
