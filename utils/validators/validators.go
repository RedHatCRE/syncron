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
package validators

import (
	"time"

	"github.com/sirupsen/logrus"
)

func ValidateTime(fromData time.Time) {
	firstAvailableDate := time.Date(2018, 1, 1, 0, 0, 0, 0, time.UTC)
	latestAvailableDate := time.Now().AddDate(0, 0, -3)
	if fromData.Before(firstAvailableDate) {
		logrus.Fatal("No available data before 2018. Please try again.")
	}
	if fromData.After(latestAvailableDate) {
		logrus.Fatal("No available data after that date. Please try again")
	}

}
