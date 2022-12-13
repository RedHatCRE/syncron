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

// The purpose of the package filter is to reduce the download
// to the specified files by filtering through the configuration's files,
// downloading only the ones that match the supplied pattern given via
// --filter flag.
package filter

import (
	"strings"

	"github.com/redhatcre/syncron/pkg/cli"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// The function Component will filter through the list of
// files supplied on the configuration files, returning a slice
// of strings that match the given pattern
func Component(filterFlag []string) []string {
	var filtered []string
	files_in_config := viper.GetStringSlice(cli.SOSReports)

	if len(filterFlag) == 0 {
		return files_in_config
	}
	for _, file := range files_in_config {
		for _, comp := range filterFlag {
			if strings.Contains(file, comp) {
				filtered = append(filtered, file)
			}
		}
	}
	if len(filtered) == 0 {
		logrus.Fatal("Filter item(s) not found")
	}
	rmDuplicates := RemoveDuplicates(filtered)
	return rmDuplicates
}

// This function iterates through the given slice and removes
// any duplicate keys in the slice, returning a list without
// any duplicates.
func RemoveDuplicates[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
