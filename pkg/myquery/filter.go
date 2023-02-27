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
	"fmt"
	"strconv"

	"github.com/redhatcre/syncron/pkg/charts"
	"github.com/redhatcre/syncron/pkg/cli"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
)

var ComponentStages map[string][]bson.E

func FilterCollections(Components []string) ([]bson.E, []bson.E) {
	filterList := []bson.E{}
	optionsList := []bson.E{}

	ComponentStages := map[string][]bson.E{
		"manila":             {bson.E{Key: "manila_conf", Value: bson.D{{Key: "$exists", Value: "true"}}}},
		"nova":               {bson.E{Key: "nova_conf", Value: bson.D{{Key: "$exists", Value: "true"}}}},
		"installed packages": {bson.E{Key: "installed_packages", Value: bson.D{{Key: "$exists", Value: "true"}}}},
		"galera":             {bson.E{Key: "galera_conf", Value: bson.D{{Key: "$exists", Value: "true"}}}},
		"ceph":               {bson.E{Key: "cephfs", Value: bson.D{{Key: "$exists", Value: "true"}}}},
		"cinder":             {bson.E{Key: "cinder_conf", Value: bson.D{{Key: "$exists", Value: "true"}}}},
		"glance":             {bson.E{Key: "glance_api_conf", Value: bson.D{{Key: "$exists", Value: "true"}}}},
		"tripleoceph":        {bson.E{Key: "cinder_conf.default.enabled_backends", Value: bson.D{{Key: "$eq", Value: "tripleo_ceph"}}}},
		"octavia-worker":     {bson.E{Key: "ps.octavia-worker", Value: bson.D{{Key: "$exists", Value: "true"}}}},
		"neutron-server":     {bson.E{Key: "ps.neutron-server", Value: bson.D{{Key: "$exists", Value: "true"}}}},
		"+3 nodes":           {bson.E{Key: "hosts_file.node_count", Value: bson.D{{Key: "$gt", Value: 3}}}},
		"none":               {bson.E{Key: "open_source_4_life", Value: bson.D{{Key: "$exists", Value: "true"}}}},
		"all":                {},
	}
	for _, component := range Components {
		c, ok := ComponentStages[component]
		if !ok {
			logrus.Fatal("No data for that component.")
		}
		if component == "all" {
			filterList = []bson.E{}
			optionsList = []bson.E{}
			break
		}
		filterList = append(filterList, c[0])
		// optionsList = append(optionsList, bson.E{Key: c[0].Key, Value: 1})
	}
	filterList = append(filterList, bson.E{
		Key: "versions_history.osp_ver", Value: bson.D{{Key: "$exists", Value: "true"}}},
	)
	return filterList, optionsList
}

func NewFilter(coll *mongo.Collection, deploymentIdPipeline primitive.A) (int32, float32) {

	cursor, err := coll.Aggregate(context.TODO(), deploymentIdPipeline)
	uniqueDep := Numbers.AllReportsWithVersion
	if err != nil {
		panic(err)
	}
	defer cursor.Close(context.TODO())
	var result bson.M
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&result)
		if err != nil {
			panic(err)
		}
	}
	numDep := (result["count"].(int32))
	percentage := float32(numDep) / float32(uniqueDep) * 100

	return numDep, percentage
}

var Data = []struct {
	OSP             string
	Deployments     int32
	Percentage      float64
	PercentageTotal float64
}{}

func FilterVersions(coll *mongo.Collection, deploymentIdPipeline primitive.A) (message string) {
	cur, err := coll.Aggregate(context.TODO(), deploymentIdPipeline)
	if err != nil {
		logrus.Fatal("Error creating cursor ", err)
	}
	defer cur.Close(context.TODO())

	Data = []struct {
		OSP             string
		Deployments     int32
		Percentage      float64
		PercentageTotal float64
	}{}
	var percentageTotal float64
	for cur.Next(context.TODO()) {
		var elem bson.M
		err := cur.Decode(&elem)
		if err != nil {
			logrus.Fatal(err)
		}
		key := fmt.Sprintf("%v", elem["_id"])
		count := elem["count"].(int32)
		percentage := float64(count) / float64(Numbers.AllReportsWithVersion) * 100
		// percentageTotal = float64(count) / float64(Numbers.AllReportsWithVersion) * 100

		Data = append(Data, struct {
			OSP             string
			Deployments     int32
			Percentage      float64
			PercentageTotal float64
		}{key, count, percentage, percentageTotal})
	}
	headers := []string{"OSP", "Deployments", "Percentage"}
	var rows [][]string
	for _, d := range Data {
		percentage := fmt.Sprintf("%.2f %%", d.Percentage)
		row := []string{d.OSP, strconv.Itoa(int(d.Deployments)), percentage}
		rows = append(rows, row)
	}
	if slices.Contains(cli.Components, "none") {
		title := fmt.Sprintf("Total deployments: %v", Numbers.AllUniqueDeployments)
		GenerateTable(title, rows, headers)
	}
	return message
}

func FilterVersionsPerComp(coll *mongo.Collection, thePipeline primitive.A) (message string) {
	var gatherVersions bson.A
	allDep := Numbers.AllReportsWithVersion
	gatherVersions = bson.A{
		//add here the stages to filter the component
		bson.M{"$match": bson.M{"versions_history.osp_ver": bson.M{"$exists": true}}},
		bson.M{
			"$project": bson.M{"firstkey": bson.M{"$arrayElemAt": []interface{}{bson.M{"$objectToArray": "$versions_history.osp_ver"}, 0}},
				"deployment_id": "$deployment_id",
			}},
		bson.M{"$group": bson.M{"_id": "$deployment_id", "uniqueValues": bson.M{"$first": "$firstkey.k"}}},
		bson.M{"$group": bson.M{"_id": "$uniqueValues", "count": bson.M{"$sum": 1}}},
	}
	gatherVersions = PrependStage(gatherVersions)

	cur, err := coll.Aggregate(context.TODO(), gatherVersions)
	if err != nil {
		logrus.Fatal("Cannot create gatherVersions cursor ", err)
	}
	defer cur.Close(context.TODO())

	data := []struct {
		OSP             string
		Deployments     int32
		Percentage      float64
		PercentageTotal float64
	}{}
	var percentageTotal float64
	for cur.Next(context.TODO()) {
		var elem bson.M
		err := cur.Decode(&elem)
		if err != nil {
			logrus.Fatal(err)
		}
		key := fmt.Sprintf("%v", elem["_id"])
		count := elem["count"].(int32)
		percentage := float64(count) / float64(UniqueDepsLen) * 100
		percentageTotal = float64(count) / float64(allDep) * 100
		data = append(data, struct {
			OSP             string
			Deployments     int32
			Percentage      float64
			PercentageTotal float64
		}{key, count, percentage, percentageTotal})
	}

	headers := []string{"OSP", "Deployments", "Percentage"}
	var rows [][]string
	for _, d := range data {
		percentage := fmt.Sprintf("%.2f %% (Total: %.2f %%)", d.Percentage, d.PercentageTotal)
		row := []string{d.OSP, strconv.Itoa(int(d.Deployments)), percentage}
		rows = append(rows, row)
	}
	if !slices.Contains(cli.Components, "none") {
		filteredTitle := fmt.Sprintf("Filtered for %v: ", cli.Components)
		GenerateTable(filteredTitle, rows, headers)
	}

	title := fmt.Sprintf("Filtered for %v", cli.Components)
	if cli.GenCharts {
		charts.CreateChart(title, "filtered.html", data)
	}

	var compM string
	for _, v := range cli.Components {
		compM += fmt.Sprintf(" %v", v)
	}

	return message
}

func PrependStage(gatherVersions primitive.A) primitive.A {
	for _, component := range cli.Components {
		var NewStage primitive.M
		switch component {
		case "cinder":
			NewStage = bson.M{"$match": bson.M{"cinder_conf": bson.M{"$exists": true}}}
		case "tripleoceph":
			NewStage = bson.M{"$match": bson.M{"cinder_conf.default.enabled_backends": "tripleo_ceph"}}
		case "none":
			return gatherVersions
		case "manila":
			NewStage = bson.M{"$match": bson.M{"manila_conf": bson.M{"$exists": true}}}
		case "nova":
			NewStage = bson.M{"$match": bson.M{"nova_conf": bson.M{"$exists": true}}}
		case "glance":
			NewStage = bson.M{"$match": bson.M{"glance_api_conf": bson.M{"$exists": true}}}
		case "octavia-worker":
			NewStage = bson.M{"$match": bson.M{"ps.octavia-worker": bson.M{"$exists": true}}}
		case "neutron-server":
			NewStage = bson.M{"$match": bson.M{"ps.neutron-server": bson.M{"$exists": true}}}
		case "+3 nodes":
			NewStage = bson.M{"$match": bson.M{"hosts_file.node_count": bson.M{"$gt": 3}}}
		}
		gatherVersions = append(bson.A{NewStage}, gatherVersions...)
	}
	return gatherVersions
}
