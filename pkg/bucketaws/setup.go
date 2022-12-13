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
	"strings"
	"time"

	"github.com/redhatcre/syncron/configuration"
	"github.com/redhatcre/syncron/pkg/parquet_reader"
	files "github.com/redhatcre/syncron/utils/files"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
)

// Builds dates from given date to current date,
// incrementing one day at a time.
func ProcessDate(fromDate time.Time) []string {

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

// Initialize a session with AWS SDK
// It will read from the file located at ~/.aws/credentials and syncron/syncron.yaml
func SetupSession(config configuration.Configuration) *session.Session {

	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(config.S3.Region),
		Endpoint: aws.String(config.S3.EndPoint),
	},
	)
	if err != nil {
		logrus.Fatal("There was an error setting up your aws session", err)
	}
	logrus.Info("Your AWS session was set up correctly")
	return sess
}

// This function initiates the service for downloading files in s3
func AccessBucket(sess *session.Session) (*s3.S3, *s3manager.Downloader) {

	svc := s3.New(sess)
	dwn := s3manager.NewDownloader(sess)

	return svc, dwn
}

// This function downloads files from bucket, filtering
// through those that are needed.
func DownloadFromBucket(config configuration.Configuration, svc *s3.S3, dwn *s3manager.Downloader, dates []string, bprefix string) error {

	var continuationToken *string
	for {
		resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
			Bucket:            aws.String(config.S3.Bucket),
			Prefix:            aws.String(files.AppendPrefix(config.Prefix, bprefix)),
			ContinuationToken: continuationToken,
		})
		if err != nil {
			logrus.Fatal("There was an error listing the objects in bucket.", err)
		}
		for _, item := range resp.Contents {
			for _, x := range dates {

				if !strings.Contains(*item.Key, x) {
					continue
				}

				logrus.Info("Downloading files for: ", bprefix)

				absoluteFileName := files.GetDownloadPath(config.DownloadDir, *item.Key)
				if files.FileExists(absoluteFileName) {
					logrus.Info("File already exists: ", absoluteFileName)
					continue
				}

				fileHandler := files.FilePathSetup(absoluteFileName)

				logrus.Info("Downloading to: ", absoluteFileName)
				start := time.Now()
				_, err := dwn.Download(
					fileHandler,
					&s3.GetObjectInput{
						Bucket: aws.String(config.S3.Bucket),
						Key:    aws.String(*item.Key),
					})
				duration := time.Since(start)
				logrus.Info("Download took: ", duration.Truncate(time.Second/2))

				if err != nil {
					fmt.Println("There was an error fetching key info.", err)
					return err
				}

				// Handle Parquet file
				noExtFileName := files.RemoveExtention(absoluteFileName)
				parquet_reader.ReadParquet(noExtFileName, absoluteFileName)
				logrus.Info("File extracted.")

				defer func() {
					if err := fileHandler.Close(); err != nil {
						logrus.Print("Error closing file handler for: ", absoluteFileName)
					}
				}()
			}
		}
		if !aws.BoolValue(resp.IsTruncated) {
			break
		}
		continuationToken = resp.NextContinuationToken
	}
	return nil
}

// CredCheck checks the credentials for a given AWS session.
// The function takes a pointer to an AWS session as its only parameter.
// It uses the Get method of the session's Credentials field to retrieve
// credentials and check if they are valid. If the Get method returns an error,
// the function logs a fatal error message. If the credentials are read successfully,
// the function logs an info message indicating that the credentials were read successfully.
// This function is useful for ensuring that the credentials used to authenticate
// the AWS session are valid before performing any operations that require authentication.
func CredCheck(sess *session.Session) {
	_, err := sess.Config.Credentials.Get()

	if err != nil {
		logrus.Fatal(
			"Error reading credentials. Check README for help.\n")
	}

	logrus.Info("Credentials read succesfully")
}
