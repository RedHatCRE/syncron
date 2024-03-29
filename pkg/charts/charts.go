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

package charts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/redhatcre/syncron/utils/files"
	"github.com/sirupsen/logrus"
)

func generatePieItems(data []struct {
	OSP             string
	Deployments     int32
	Percentage      float64
	PercentageTotal float64
}) []opts.PieData {
	items := make([]opts.PieData, 0)
	for _, d := range data {
		items = append(items, opts.PieData{
			Name:  fmt.Sprintf("OSP %v -> Deployments: ", d.OSP),
			Value: d.Deployments,
		})
	}
	return items
}

func CreateChart(title string, chartPath string, data []struct {
	OSP             string
	Deployments     int32
	Percentage      float64
	PercentageTotal float64
}) {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeChalk}),
		charts.WithTitleOpts(opts.Title{
			Title: title,
		}))
	pie.SetSeriesOptions()
	pie.AddSeries("Versions",
		generatePieItems(data)).
		SetSeriesOptions(
			charts.WithPieChartOpts(
				opts.PieChart{
					Radius: 200,
				},
			),
			charts.WithLabelOpts(
				opts.Label{
					Show:      true,
					Formatter: "{b}: {c}",
				},
			),
		)
	files.FilePathSetup("/tmp/syncron/charts/" + chartPath)
	f, _ := os.Create(filepath.Clean("/tmp/syncron/charts/" + chartPath))
	defer func() {
		if err := f.Close(); err != nil {
			logrus.Errorf("Error closing file: %s\n", err)
		}
	}()
	err := pie.Render(f)
	if err != nil {
		logrus.Error("Error rendering chart \n", err)
	}
}

func CreateChartBars(chartPath string, numbers []int32, components []string) {

	components = trimSlice(components)
	var barData []opts.BarData
	for _, d := range numbers {
		barData = append(barData, opts.BarData{Value: d})
	}
	bar := charts.NewBar()

	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Theme: types.ThemeChalk}),
		charts.WithTitleOpts(opts.Title{Title: "Deployments per feature"}))
	bar.SetXAxis(components).AddSeries("deployments", barData)
	files.FilePathSetup("/tmp/syncron/charts/" + chartPath)
	f, _ := os.Create(filepath.Clean("/tmp/syncron/charts/" + chartPath))
	defer func() {
		if err := f.Close(); err != nil {
			logrus.Errorf("Error closing file: %s\n", err)
		}
	}()
	err := bar.Render(f)
	if err != nil {
		logrus.Error("Error rendering chart \n", err)
	}
}

func trimSlice(s []string) []string {
	result := make([]string, len(s))
	for i, str := range s {
		if len(str) > 3 {
			result[i] = str[0:3]
		}
	}
	return result
}
