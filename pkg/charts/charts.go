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

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"github.com/sirupsen/logrus"
)

func generatePieItems(data []struct {
	OSP         string
	Deployments int32
	Percentage  float64
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
	OSP         string
	Deployments int32
	Percentage  float64
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
	f, _ := os.Create(chartPath)
	defer f.Close()
	err := pie.Render(f)
	if err != nil {
		logrus.Error("Error rendering chart \n", err)
	} else {
		logrus.Infof("Generated chart at %s", chartPath)
	}

}
