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
	"time"

	s3setup "github.com/redhatcre/syncron/pkg/bucketaws"
	"github.com/redhatcre/syncron/pkg/cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var download = &cobra.Command{
	Use:       "download",
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	ValidArgs: []string{cli.Insights, cli.SOSReports, cli.All},
	RunE:      onRun,
	Short:     "Download files from bucket",
}

func init() {

	download.Flags().Int(
		cli.Days,
		2,
		"Download data from the last X days",
	)
	download.Flags().Int(
		cli.Months,
		0,
		"Download data from the last X months",
	)

	download.Flags().Int(
		cli.Years,
		0,
		"Download data from the last X years",
	)
	root.AddCommand(download)
}

func onRun(cmd *cobra.Command, args []string) error {
	Month, _ := cmd.Flags().GetInt(cli.Months)
	Year, _ := cmd.Flags().GetInt(cli.Years)
	Day, err := cmd.Flags().GetInt(cli.Days)

	if Day < 2 {
		logrus.Error("No data available. Try again with 3 or more days.")
		return err
	}
	fromDate := time.Now().AddDate(-Year, -Month, -Day)

	// Reading configuration file
	s3setup.ConfigRead()
	// Creating AWS session
	sess := s3setup.SetupSession()
	//Checking credentials
	s3setup.Credcheck(sess)
	// Processing dates to download
	dates := s3setup.ProcessDate(fromDate)
	// Accessing bucket
	svc, dwn := s3setup.AccessBucket(sess)

	logrus.Info("Pulling data since ",
		fromDate.Year(),
		fromDate.Month(),
		fromDate.Day())

	filesToDownload := viper.GetStringSlice(cli.SOSReports)
	for _, f := range filesToDownload {
		logrus.Info("Downloading files for: ", f)
		s3setup.DownloadFromBucket(svc, dwn, dates, f)
	}
	return nil
}
