package main

import (
	"log"

	"github.com/rhcre/syncron/cmd/adhoc/cmd"
)

func main() {
	err := cmd.Execute()

	if err != nil {
		log.Fatal(err)
	}
}
