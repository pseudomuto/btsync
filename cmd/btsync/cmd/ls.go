package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"context"
	"fmt"
	"os"

	"github.com/pseudomuto/btsync/pkg/schema"
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list partition details",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		if err := listPartitionDirectory(ctx, viper.GetString("dir")); err != nil {
			cmd.Println(err)
			os.Exit(1)
		}
	},
}

func listPartitionDirectory(ctx context.Context, dir string) error {
	tw := NewTableWriter(os.Stdout)
	tw.SetHeader([]string{"Partition", "Name", "Description"})

	for res := range schema.ParseDirectory(ctx, dir) {
		if res.Err != nil {
			return res.Err
		}

		for _, table := range res.Partition.Tables {
			tw.Append([]string{res.Partition.Name, table.Name, table.Description})
		}
	}

	if tw.NumLines() == 0 {
		return fmt.Errorf("no partitions found in: %s", dir)
	}

	tw.Render()
	return nil
}
