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

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/redhatcre/syncron/pkg/charts"
	"github.com/redhatcre/syncron/pkg/cli"
	"github.com/redhatcre/syncron/pkg/myquery"
	pipes "github.com/redhatcre/syncron/pkg/pipelines"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var query = &cobra.Command{
	Use:       "queries",
	Short:     "Query database",
	Args:      cobra.MatchAll(cobra.OnlyValidArgs),
	ValidArgs: []string{cli.AccountN, cli.QueryID},
	RunE:      onRunQuery,
}

func init() {
	query.Flags().String(
		cli.AccountN,
		"",
		"Specify account number to pull data from.",
	)
	query.Flags().String(
		cli.QueryID,
		"",
		"Specify query ID to pull data from.",
	)
	query.Flags().String(
		cli.Output,
		"",
		"Choose output for report ",
	)
	query.Flags().Bool(
		cli.GenerateCharts,
		false,
		"Choose output for report ",
	)
	root.AddCommand(query)
}

func onRunQuery(cmd *cobra.Command, args []string) error {

	myquery.GetUserInfo()
	m := myquery.MyDataBase{}
	m.CreateClient()
	defer disconnectClient(m)

	// This is our main loop
	for {
		cli.MainSelector()
		if cli.DeploymentID != "" {
			OneDeployment(m)
			continue
		}
		// Are we going to draw charts?
		cli.GenCharts = cli.ChartSelector(cli.Components)
		// Are we going to get one deployment?
		filterList, optionsList := myquery.FilterCollections(cli.Components)
		if !slices.Contains(cli.Components, "none") {
			logrus.Infof("Features: %v", cli.Components)
		}
		opts := options.Find().SetProjection(optionsList)
		filterComponents := bson.D(filterList)
		if cli.Components[0] != "none" {
			myquery.PrintResults(m.Coll, filterComponents, opts)
		}
		countDeployments(m.Coll, filterComponents, m.MyClient)
		if cli.GenCharts {
			charts.CreateChart("Total deployments stats", "total.html", myquery.Data)
			logrus.Info("Charts created at /tmp/syncron/charts/")
		}
	}
}

func disconnectClient(m myquery.MyDataBase) {
	if err := m.MyClient.Disconnect(context.TODO()); err != nil {
		logrus.Fatal("Could not disconnect", err)

	}
}

func OneDeployment(m myquery.MyDataBase) {
	if cli.DeploymentID != " " {
		var v myquery.Deployment
		var result bson.M
		filter := bson.D{{Key: "deployment_id", Value: cli.DeploymentID}}
		errFindOne := m.Coll.FindOne(context.TODO(), filter).Decode(&result)
		if errFindOne != nil {
			logrus.Error(errFindOne)
		} else {
			output, err := json.MarshalIndent(result, "", "   ")

			if err := json.Unmarshal(output, &v); err != nil {
				logrus.Fatal("Error unmarshaling output", err)
			}
			if err != nil {
				if err == mongo.ErrNoDocuments {
					logrus.Fatal("No docs found", err)
				}
				panic(err)
			}
			CinderDrivers := strings.Join(v.CinderConf.Driver, " ")
			myquery.DepInfo = map[string]string{
				"Acc.Number":                  v.AccountNumber,
				"DB-MaxEntries":               v.CinderConf.DatabaseKey.MaxRetries,
				"EnabledV3":                   v.CinderConf.Default.EnabledV3api,
				"TripleoNFS":                  strconv.FormatBool(v.CinderConf.TripleoNFS.NasSecureFileOperations),
				"Cinder Drivers":              CinderDrivers,
				"OSP Api Workers":             v.CinderConf.Default.OSApiVolumeWorkers,
				"Available Cinder Ztore zone": v.CinderConf.Default.StoraAvailZone,
				"Key Manager":                 v.CinderConf.KeyManager.Backend,
			}
			title := fmt.Sprintf("Deployment info for ID: %v", cli.DeploymentID)
			var rows [][]string
			for k, v := range myquery.DepInfo {
				rows = append(rows, []string{k, v})
			}
			headers := []string{"Key", "Value"}
			myquery.GenerateTable(title, rows, headers)
		}
	}
}

func countDeployments(coll *mongo.Collection, filterComponents primitive.D, client *mongo.Client) {

	s := cli.IniSpinner(" Fetching data from database")

	pipes.StoragePipes = pipes.InitStoragePipes()
	pipes.CommonPipes = pipes.InitCommonPipes()
	pStorage := pipes.StoragePipes
	pCommon := pipes.CommonPipes

	filter := myquery.NewFilter
	var message string
	allReports, err := coll.EstimatedDocumentCount(context.TODO())
	myquery.Numbers.AllReports = allReports
	if err != nil {
		logrus.Fatal("Error counting all documents.", err)
	}

	uniqueDeploymentsNumber, _ := filter(coll, pCommon.UniqueDeployments)
	totalWithVersionR, _ := filter(coll, pCommon.TotalWithVersion)

	myquery.Numbers.AllUniqueDeployments = uniqueDeploymentsNumber
	myquery.Numbers.AllReportsWithVersion = totalWithVersionR

	myquery.Manila.Deployments, myquery.Manila.Percentage = filter(coll, pStorage.ManilaEnabled)
	cephEnabledResults, cephEnabledPercentage := filter(coll, pStorage.CephEnabled)
	glance_exists, glancePercentage := filter(coll, pStorage.GlanceConf)
	nova_exists, novaPercentage := filter(coll, pStorage.NovaConf)
	cephfs_exists, cephfsPercentage := filter(coll, pStorage.Cephfs_enabled)

	hciDeployR, hciPercentage := filter(coll, pStorage.HciDeploy)
	nfsManilaR, nfsManilaPercentage := filter(coll, pStorage.NFSManila)
	directorDeployR, directorPercentage := filter(coll, pStorage.DirectorDeployed)
	tripleoCephR, tripleoCephPercentage := filter(coll, pStorage.TripleoCeph)
	listDatabases(message, client)

	// Stop spinner
	s.Stop()
	if slices.Contains(cli.Components, "none") {
		logrus.Infof("Total number of SOSreports: %d\n", myquery.Numbers.AllReports)
	}
	myquery.FilterVersions(coll, pCommon.GatherVersions)
	var components []string
	var depNumbers []int32
	var data = []struct {
		component   string
		deployments interface{}
		percentage  float32
	}{
		{"Glance", glance_exists, glancePercentage},
		{"Nova", nova_exists, novaPercentage},
		{"Ceph", cephEnabledResults, cephEnabledPercentage},
		{"Manila", myquery.Manila.Deployments, myquery.Manila.Percentage},
		{"Ceph-fs", cephfs_exists, cephfsPercentage},
		{"Director deployed", directorDeployR, directorPercentage},
		{"Tripleo managed ceph", tripleoCephR, tripleoCephPercentage},
		{"HCI deploy", hciDeployR, hciPercentage},
		{"NFS deploy", nfsManilaR, nfsManilaPercentage},
	}

	headers := []string{"OSP", "Deployments", "Percentage"}
	var rows [][]string
	var row []string
	for _, d := range data {
		percentage := fmt.Sprintf("%.2f %%", d.percentage)
		row = []string{d.component, fmt.Sprintf("%v", d.deployments), percentage}
		depNumbers = append(depNumbers, d.deployments.(int32))
		components = append(components, row[0])
		rows = append(rows, row)
	}
	myquery.FilterVersionsPerComp(coll, pStorage.NFSManila)
	if slices.Contains(cli.Components, "none") {
		myquery.GenerateTable("Per component", rows, headers)
		charts.CreateChartBars("percomponent.html", depNumbers, components)
	}

}

func listDatabases(message string, client *mongo.Client) {
	tab := "- "
	message += "Available databases: \n"
	databases, _ := client.ListDatabaseNames(context.Background(), bson.D{})
	for _, x := range databases {
		message += fmt.Sprintf("%s%s\n", tab, x)
	}
}
