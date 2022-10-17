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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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
	resp, err := svc.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: aws.String(viper.GetString("bucket"))})
	if err != nil {
		fmt.Println("Unable to list items in bucket", resp)
		return err
	} else {
		logrus.Info("Success getting into the bucket!")
	}
	return nil
}

func DownloadFromBucket() error {

	// Future PR implementation

	return nil
}
