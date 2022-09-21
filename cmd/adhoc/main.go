package main

import (
	cmd "github.com/rhcre/syncron/cmd/adhoc/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	err := cmd.Execute()

	if err != nil {
		log.Fatal(err)
	}
}
