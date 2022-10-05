package main

import (
	"github.com/redhatcre/syncron/cmd/adhoc/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	err := cmd.Execute()

	if err != nil {
		logrus.Fatal(err)
	}
}
