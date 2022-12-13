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

package parquet_reader

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/apache/arrow/go/v11/parquet/file"

	"github.com/redhatcre/syncron/pkg/cli"
	"github.com/redhatcre/syncron/pkg/dumper"
)

// This function reads a Parquet file and prints its contents to an output file or to standard output.
// The function accepts the name of the output file as a parameter, as well as the name of the file to be read.
// The function has a config struct that allows the user to specify whether to print key-value metadata and whether
// to use a memory map when reading the file. The function uses the file package to open the Parquet file,
// reads its metadata, and creates a dumper for each column specified by the user.
// The function then loops through each row group in the file,
// printing the selected columns and their values to the output file or standard output.
func ReadParquet(outputFile string, fileToRead string) {
	var config struct {
		NoMemoryMap bool
		//Columns     string
	}
	columns := "2"

	var dataOut io.Writer
	dataOut = os.Stdout
	if outputFile != "-" {
		var err error
		fileOut, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: output %q cannot be created, %s\n", cli.Output, err)
			os.Exit(1)
		}
		bufOut := bufio.NewWriter(fileOut)
		defer func() {
			bufOut.Flush()
			fileOut.Close()
		}()
		dataOut = bufOut
	}
	// This if statement checks if the Columns field in the config variable is not an empty string.
	// If it isn't, the statement splits the Columns field on commas and then iterates over the resulting list of strings.
	// For each string, it attempts to convert it to an integer using the strconv.Atoi function.
	// If this conversion succeeds, the integer value is appended to the selectedColumns slice.
	selectedColumns := []int{}
	if columns != "" {
		for _, c := range strings.Split(columns, ",") {
			cval, err := strconv.Atoi(c)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error: selectedColumns needs to be comma-delimited integers")
				os.Exit(1)
			}
			selectedColumns = append(selectedColumns, cval)
			//uniqueValues := filter.RemoveDuplicates(selectedColumns)
			fmt.Println("Original path to this file is: ", selectedColumns)
		}
	}

	rdr, err := file.OpenParquetFile(fileToRead, !config.NoMemoryMap)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error opening parquet file: ", err)
		os.Exit(1)
	}
	
	fileMetadata := rdr.MetaData()

	// If statement below checks whether the selectedColumns slice is empty. If it is empty,
	// it appends the index of each column in the schema to the selectedColumns slice.
	// Otherwise, it iterates over each column in the selectedColumns slice
	// and checks if its index is out of range. If it is out of range,
	//it prints an error message and exits with an error code of 1.
	if len(selectedColumns) == 0 {
		for i := 0; i < fileMetadata.Schema.NumColumns(); i++ {
			selectedColumns = append(selectedColumns, i)
		}
	} else {
		for _, c := range selectedColumns {
			if c < 0 || c >= fileMetadata.Schema.NumColumns() {
				fmt.Fprintln(os.Stderr, "selected column is out of range")
				os.Exit(1)
			}
		}
	}

	for r := 0; r < rdr.NumRowGroups(); r++ {
		rgr := rdr.RowGroup(r)
		const colwidth = 18

		scanners := make([]*dumper.Dumper, len(selectedColumns))
		for idx, c := range selectedColumns {
			col, err := rgr.Column(c)
			if err != nil {
				log.Fatalf("unable to fetch column=%d err=%s", c, err)
			}
			scanners[idx] = dumper.CreateDumper(col)
			fmt.Fprintf(dataOut, fmt.Sprintf("%%-%ds|", colwidth), col.Descriptor().Name())
		}
		fmt.Fprintln(dataOut)

		// This for loop iterates over the scanners and fetches the next value
		// from each scanner using the Next method. If the value exists, it is
		// formatted using the FormatValue method and printed to the dataOut
		// writer along with a "|" character. If the value does not exist, the
		// loop adds an empty string of the appropriate width to the line variable
		// or prints an empty string to the dataOut writer, depending on whether
		// data has been printed in the current iteration. The loop breaks if no
		// data is printed in an iteration.
		//var line string
		for {
			data := false
			for _, s := range scanners {
				if val, ok := s.Next(); ok {
					//if !data {
					//	fmt.Fprint(dataOut, line)
					//}
					fmt.Fprint(dataOut, s.FormatValue(val, colwidth), "|")
					data = true
				} //else {
				//	if data {
				//		fmt.Fprintf(dataOut, fmt.Sprintf("%%-%ds|", colwidth), "")
				//	} else {
				//		line += fmt.Sprintf(fmt.Sprintf("%%-%ds|", colwidth), "")
				//	}
				//}//
			}
			if !data {
				break
			}
			//fmt.Fprintln(dataOut)
			//line = ""
		}
		//_, err := fmt.Fprintln(dataOut)
		//if err != nil {
		//	logrus.Error("Following error reading line", err)
		//}
	}
}
