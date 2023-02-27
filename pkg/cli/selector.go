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

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
)

var GenCharts bool
var Components []string
var OneDepByID bool

func MainSelector() {
	selection := " "
	DeploymentID = ""
	Components = []string{"none"}
	prompt := &survey.Select{
		Message: "Choose what you need: ",
		Options: []string{
			"Filter Deployments",
			"Get deployment by deploymentID",
			"Just let me browse",
			"Quit",
		},
	}

	if err := survey.AskOne(prompt, &selection); err != nil {
		logrus.Error("Error on main selector", err)
	}
	switch {
	case selection == "Filter Deployments":
		Components = ComponentSelector()
	case selection == "Get deployment by deploymentID":
		DeploymentID = DeploymentIDSelector()
	case selection == "Quit":
		os.Exit(0)
	}
}
func DeploymentIDSelector() string {
	depID := " "
	for {
		prompt := &survey.Input{
			Message: "Input the deployment ID: ",
		}
		if err := survey.AskOne(prompt, &depID); err != nil {
			logrus.Error("Error getting deployment ID", err)
		}
		if depID == "" {
			logrus.Warn("Sorry, deploymentID cannot be empty")
			continue
		}
		return depID
	}
}

func ComponentSelector() []string {
	selection := []string{}
	prompt := &survey.MultiSelect{
		Message: "What filter to apply: ",
		Options: []string{
			"manila",
			"cinder",
			"tripleoceph",
			"nova",
			"glance",
			"octavia-worker",
			"neutron-server",
			"+3 nodes",
			"quit"},
	}

	if err := survey.AskOne(prompt, &selection); err != nil {
		logrus.Error("Error filtering features", err)
	}
	if slices.Contains(selection, "quit") {
		fmt.Print("\033[1A\033[K")
		os.Exit(0)
	}
	if len(selection) == 0 {
		return []string{"none"}
	}
	return selection
}

func ChartSelector(selection []string) bool {
	if len(selection) == 0 {
		return false
	}
	this := false

	prompt := &survey.Confirm{
		Message: "Would you like to generate charts?",
	}
	if err := survey.AskOne(prompt, &this); err != nil {
		logrus.Error("Error selecting if generating chart", err)
	}
	return this
}
