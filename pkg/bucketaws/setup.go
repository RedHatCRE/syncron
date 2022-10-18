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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ConfigRead() error {
	// Setting up file formatting
	// Pulling from Viper
	viper.SetConfigFile("config/syncron.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	// Reading from file
	err := viper.ReadInConfig()
	if err != nil {
		return err
	} else {
		logrus.Info("Your configuration file was read succesfully")
		logrus.Info("Reading from bucket: ", viper.Get("bucket"))
	}
	return nil
}

func SetupSession() *session.Session {
	// Initialize a session with AWS SDK
	// It will read from the file located at ~/.aws/credentials and syncron/config/syncron.yaml
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(viper.GetString("s3.region")),
		Endpoint: aws.String(viper.GetString("s3.endpoint")),
	},
	)
	if err != nil {
		fmt.Println("There was an error setting up your aws session")
		os.Exit(1)
	} else {
		logrus.Info("Your AWS session was set up correctly")
	}
	return sess
}

func AccessBucket(sess *session.Session) error {
	logrus.Info("Accessing bucket...")
	svc := s3.New(sess)
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(viper.GetString("bucket")),
		Prefix: aws.String("attachment_data/sosreport/sos_extraction_rules.sos.sos_run_info/created_year=2022/created_month=10/created_day=5/"),
	})
	for _, item := range resp.Contents {
		DownloadFromBucket(item.Key)
		fmt.Println(*item.Key, item.LastModified)
		//	fmt.Println("Size:         ", *item.Size)
		//	fmt.Println("Storage class:", *item.StorageClass)
		//	fmt.Println("")
	}
	fmt.Printf("Number of items: %d\n", len(resp.Contents))

	if err != nil {
		fmt.Println("Unable to list items in bucket", resp)
		return err
	} else {
		logrus.Info("Success getting into the bucket!")
	}
	return nil
}

func DownloadFromBucket(key *string) error {

	sess := SetupSession()

	dwn := s3manager.NewDownloader(sess)
	filePath := strings.Split(*key, "/")
	os.MkdirAll("/tmp/syncron/" + *key, 0700)
	fooFile, err := os.Create("/tmp/" + filePath[len(filePath)-1])
	if err != nil {
		fmt.Println(err)
	} else {
		logrus.Info("File opened correctly")
	}
	objects := []s3manager.BatchDownloadObject{
		{
			Object: &s3.GetObjectInput{
				Bucket: aws.String("DH-STAGE-INSIGHTS"),
				Key:    aws.String(*key),
			},
			Writer: fooFile,
		},
	}

	iter := &s3manager.DownloadObjectsIterator{Objects: objects}

	if err := dwn.DownloadWithIterator(aws.BackgroundContext(), iter); err != nil {
		fmt.Println("nooooo", err)
	}

	return nil
}
