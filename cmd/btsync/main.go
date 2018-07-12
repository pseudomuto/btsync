package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pseudomuto/btsync/pkg/config"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	for res := range config.ParseDirectory(ctx, "testdata/partitions") {
		if res.Err != nil {
			fmt.Printf("ERROR: %s\n", res.Err.Error())
			continue
		}

		fmt.Printf("Partition: %+v\n", res.Partition)
		for _, tbl := range res.Partition.Tables {
			fmt.Printf("Table: %+v\n", tbl)
		}
	}

	fmt.Println("Done")
}
