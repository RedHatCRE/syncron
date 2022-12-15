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

package snappy_reader

import (
	"os"

	"github.com/golang/snappy"
	"github.com/sirupsen/logrus"
)

// CompressFile compresses the given file using the snappy algorithm
func CompressFile(inputPath string, outputPath string) error {
	// Read the input file
	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// Compress the input data
	compressedData := snappy.Encode(nil, inputData)

	// Write the compressed data to the output file
	err = os.WriteFile(outputPath, compressedData, 0644)
	if err != nil {
		return err
	}

	return nil
}

// DecompressFile decompresses the given file using the snappy algorithm
func DecompressFile(inputPath string, outputPath string) error {
	// Read the input file
	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		logrus.Error("Error reading file")
		return err
	}

	// Decompress the input data
	decompressedData, err := snappy.Decode(nil, inputData)
	if err != nil {
		logrus.Error("Error decoding data")
		return err
	}

	// Write the decompressed data to the output file
	err = os.WriteFile(outputPath, decompressedData, 0644)
	if err != nil {
		logrus.Error("Error writing file")
		return err
	}

	return nil
}
