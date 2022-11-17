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
package cli

// Keywords used as arguments on the CLI.
const (
	Insights   string = "insights"
	SOSReports string = "sosreports"
	All        string = "all"
)

// Keywords used as flags on the CLI.
const (
	Debug string = "debug"

	Days   string = "days"
	Months string = "months"
	Years  string = "years"
	Filter string = "filter"
)

// Keywords used as shortcuts on the CLI.
const (
	D string = "d"
)
