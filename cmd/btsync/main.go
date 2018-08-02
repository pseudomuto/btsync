package main

import (
	"log"

	"github.com/pseudomuto/btsync/cmd/btsync/cmd"
)

func main() {
	if err := cmd.Execute(Configure); err != nil {
		log.Fatal(err)
	}
}
