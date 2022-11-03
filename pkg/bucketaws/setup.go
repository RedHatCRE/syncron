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

package s3setup

import (
	"fmt"
	"os"
	"strings"
	"time"

	files "github.com/redhatcre/syncron/utils"


	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ProcessDate(fromDate time.Time) []string {

	// This function formats all dates from provided date
	// until current date and returns a string formatted
	// for s3 keys.

	start := fromDate
	end := time.Now()
	var dates []string
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		year, month, day := d.Date()
		dates = append(
			dates,
			fmt.Sprintf(
				"created_year=%d/created_month=%d/created_day=%d",
				year, month, day))
	}
	return dates
}

func ConfigRead() error {

	// Setting up file formatting
	// Using Viper
	// Reading from file syncron.yaml

	viper.AddConfigPath("./config")
	viper.SetConfigFile("syncron.yaml")
	viper.SetConfigType("yaml")
	viper.SetDefault("download_dir", "/tmp/syncron/")

	// Reading from file
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Error(err)
		logrus.Exit(1)
	}
	logrus.Info("Your configuration file was read succesfully")
	logrus.Info("Reading from bucket: ", viper.Get("bucket"))
	return nil
}

func SetupSession() *session.Session {

	// Initialize a session with AWS SDK
	// It will read from the file located at ~/.aws/credentials and syncron/syncron.yaml

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(viper.GetString("s3.region")),
		Endpoint: aws.String(viper.GetString("s3.endpoint")),
	},
	)
	if err != nil {
		fmt.Println("There was an error setting up your aws session", err)
		os.Exit(1)
	}
	logrus.Info("Your AWS session was set up correctly")
	return sess
}

func AccessBucket(sess *session.Session) (*s3.S3, *s3manager.Downloader) {

	// This function initiates the service for downloading files in s3

	logrus.Info("Accessing bucket...")
	svc := s3.New(sess)
	dwn := s3manager.NewDownloader(sess)

	return svc, dwn
}

func DownloadFromBucket(svc *s3.S3, dwn *s3manager.Downloader, dates []string, bprefix string) error {

	// This function takes care of listing the keys in the bucket, filtering
	// through those that are needed.

	var continuationToken *string

	for {
		resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:            aws.String(viper.GetString("bucket")),
			Prefix:            aws.String(files.AppendPrefix(bprefix)),
			ContinuationToken: continuationToken,
		})
		if err != nil {
			logrus.Error("There was an error listing the objects in bucket.")
			logrus.Error(err)
			os.Exit(1)
		}
		for _, item := range resp.Contents {
			for _, x := range dates {
				if strings.Contains(*item.Key, x) {
					fooFile, fileName := files.FilePathSetup(item.Key, dwn)
					logrus.Info("Downloading ", fileName)
					_, err := dwn.Download(
						fooFile,
						&s3.GetObjectInput{
							Bucket: aws.String(viper.GetString("bucket")),
							Key:    aws.String(*item.Key),
						})
					if err != nil {
						fmt.Println("There was an error fetching key info.", err)
						return err
					}
				}
			}
		}
		if !aws.BoolValue(resp.IsTruncated) {
			break
		}
		continuationToken = resp.NextContinuationToken
	}
	return nil
}
