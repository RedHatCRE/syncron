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
	"path"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// Creates a folder/file given the path in a recursive way
// If the path finish with '/' it will create a folder
// Otherwise it will create a file
func FilePathSetup(absoluteFilePath string) *os.File {
	err := os.MkdirAll(filepath.Dir(absoluteFilePath), 0750)

	if err != nil {
		logrus.Fatal(err)
	}

	fileHandler, err := os.Create(filepath.Clean(absoluteFilePath))

	if err != nil {
		logrus.Fatal(err)
	}

	return fileHandler
}

// Check if file exists
func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}

// This function appends the prefix to components
// Make sure prefix variable is initiated  properly in config file
func AppendPrefix(prefix string, bPrefix string) string {
	return fmt.Sprint(prefix + bPrefix)
}

// Get the absolute download path of the file based on the folder and the file
func GetDownloadPath(downloadDir string, destinationFile string) string {
	return filepath.Join(filepath.Clean(downloadDir), filepath.Clean(destinationFile))
}

// This function is used to remove the file extension from a file name.
// It takes in a string representing the file name and returns a new string
//with the file extension removed.
func RemoveExtension(fileName string) string {
	fileExt := path.Ext(fileName)
	return strings.TrimSuffix(fileName, fileExt)
}
