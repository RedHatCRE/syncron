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

package myquery

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/redhatcre/syncron/pkg/cli"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slices"
)

var Deployments []Deployment
var DepInfo map[string]string
var UniqueDepsLen int

func isVersionPresent(dep Deployment) bool {
	v := dep.Versions.OSPv
	a := reflect.ValueOf(v)
	var isVersionPresent bool
	isVersionPresent = false
	for i := 0; i < a.NumField(); i++ {
		set := a.Field(i).Interface()
		if set != "" {
			isVersionPresent = true
		}
	}
	return isVersionPresent
}

func groupUniqueDep(deps []Deployment) map[string][]Deployment {
	groups := make(map[string][]Deployment)
	for _, v := range deps {
		if isVersionPresent(v) {
			groups[v.DeploymentID] = append(groups[v.DeploymentID], v)
		}
	}
	return groups
}

var cinderBackends []string
var glanceDrivers []string

func getCinderDrivers(deps []Deployment) map[string]map[string][]Deployment {
	depsByID := make(map[string][]Deployment)
	for _, v := range deps {
		depsByID[v.DeploymentID] = append(depsByID[v.DeploymentID], v)
	}

	groups := make(map[string]map[string][]Deployment)
	for _, deployments := range depsByID {
		version := GetVersion(deployments[0])

		if _, ok := groups[version]; !ok {
			groups[version] = make(map[string][]Deployment)
		}
		groups[version][deployments[0].DeploymentID] = deployments
	}

	for _, deployments := range groups {
		cinderBackends = getCinderBackends(deployments)
		glanceDrivers = getGlanceDrivers(deployments)
	}
	if slices.Contains(cli.Components, "glance") {
		logrus.Infof("Number of glance drivers found %v", len(glanceDrivers))
	}
	if slices.Contains(cli.Components, "cinder") {
		logrus.Infof("Number of Cinder Backends found %v", len(cinderBackends))
	}

	return groups
}

func getCinderBackends(dep map[string][]Deployment) []string {
	var backends []string
	for _, v := range dep {
		for _, v := range v {
			b := v.CinderConf.Default.EnabledBackends
			splitbackends := strings.Split(b, ",")
			for _, v := range splitbackends {
				if !slices.Contains(backends, v) && v != "" {
					backends = append(backends, v)
				}
			}
		}
	}
	return backends
}

func getGlanceDrivers(dep map[string][]Deployment) []string {
	var Drivers []string
	var unDriver []string
	var Driver string
	driverCount := make(map[string]int)

	for _, v := range dep {
		for _, v := range v {
			b := v.GlanceConf.Default.GlanceDrivers
			if b != "" {
				splitDrivers := strings.Split(b, ":")
				if splitDrivers[0] == "default_backend" {
					unDriver = strings.Split(splitDrivers[1], ",")
					Driver = unDriver[0]
					if !slices.Contains(Drivers, Driver) {
						Drivers = append(Drivers, Driver)
					}
					driverCount[Driver]++
				}
			}
		}
	}
	if slices.Contains(cli.Components, "glance") {
		for driver, count := range driverCount {
			fmt.Printf("Driver: %s, Count: %d\n", driver, count)
		}
	}
	return Drivers
}

func GetVersion(dep Deployment) string {
	v := dep.Versions.OSPv
	a := reflect.ValueOf(v)
	typeOfS := a.Type()
	var version string
	for i := 0; i < a.NumField(); i++ {
		Vstr := typeOfS.Field(i).Name
		version = Vstr
	}
	return version
}

func PrintResults(coll *mongo.Collection, filterComponents primitive.D, opts *options.FindOptions) {
	cursor, err := coll.Find(context.TODO(), filterComponents, opts)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO())
	c := cli.IniSpinner(" Fetching data from database...")
	for cursor.Next(context.TODO()) {
		var result bson.M
		if err = cursor.Decode(&result); err != nil {
			panic(err)
		}

		output, err := json.MarshalIndent(result, "", "   ")
		if err != nil {
			panic(err)
		}
		var dep Deployment
		errJson := json.Unmarshal(output, &dep)
		if errJson != nil {
			logrus.Error("Error unmarshaling\n", errJson)
		}
		Deployments = append(Deployments, dep)
	}
	c.Stop()
	uniqueDeps := groupUniqueDep(Deployments)
	if slices.Contains(cli.Components, "cinder") {
		getCinderDrivers(Deployments)
	}
	UniqueDepsLen = len(uniqueDeps)
	logrus.Infof("Your query contains %v deployments and %v are unique", len(Deployments), len(uniqueDeps))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			logrus.Fatal("Couldn't get any documents.\n", err)
		}
	}
	Deployments = []Deployment{}
}

func GenerateTable(title string, rows [][]string, headers []string) string {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.Style().Title.Align = text.AlignCenter

	t.SetTitle(title)
	headerRow := table.Row{}
	for _, header := range headers {
		headerRow = append(headerRow, header)
	}

	t.AppendHeader(headerRow)

	for _, row := range rows {
		tableRow := table.Row{}
		for _, cell := range row {
			tableRow = append(tableRow, cell)

		}
		t.AppendRow(tableRow)
	}

	t.SetStyle(table.StyleRounded)

	return t.Render()
}
