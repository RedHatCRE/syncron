// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This is a reimplementation of the amazing work done by Mathew Topol.
// Thanks to his work, reading parquet files was made incredibly easy.

package parseador

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// ParseLines reads the contents of a file and returns the lines that contain a certain element
func ParseLines(file string, element string) ([]string, error) {
	var lines []string

	// Open the file
	f, err := os.Open(filepath.Clean(file))
	if err != nil {
		return lines, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			logrus.Printf("Error closing file: %s\n", err)
		}
	}()

	// Read the file line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line contains the element
		if strings.Contains(line, element) {
			if !strings.HasPrefix(line, "#") {
				lines = append(lines, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return lines, err
	}

	return lines, nil
}

func OutputData(file string, data []string) error {
	// Marshal the data into JSON
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Write the JSON to the file
	err = os.WriteFile(file, bytes, 0600)
	if err != nil {
		return err
	}

	return nil
}
