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

package files

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// This function takes a key from s3 and prepares the system for
// download by recursively creating the path needed.
func FilePathSetup(key *string, dwn *s3manager.Downloader) (*os.File, string) {

	download_dir := viper.GetString("download_dir")
	err := os.MkdirAll(download_dir+filepath.Dir(*key), 0700)
	if err != nil {
		logrus.Fatal("An error ocurred creating paths", err)
	}
	fileName := filepath.Base(*key)

	fileHandler, err := os.Create(filepath.Clean(download_dir + filepath.Dir(*key) + "/" + fileName))

	if err != nil {
		fmt.Println(err)
	}

	return fileHandler, fileName
}

// This function appends the prefix to components
// Make sure prefix variable is initiated  properly in config file
func AppendPrefix(bPrefix string) string {

	prefix := viper.GetString("prefix")
	fullPrefix := fmt.Sprint(prefix + bPrefix)
	return fullPrefix
}

func GetConfigPath() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		logrus.Fatal(err)
	}
	return dirname + "/.config/"
}
