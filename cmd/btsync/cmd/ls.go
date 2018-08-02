package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pseudomuto/btsync/pkg/config"
)

func init() {
	lsCmd.Flags().Bool("remote", false, "list BT partitions")
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
	tw.SetHeader([]string{"Name", "Description", "Tables"})

	for res := range config.ParseDirectory(ctx, dir) {
		if res.Err != nil {
			return res.Err
		}

		tables := make([]string, len(res.Partition.Tables))
		for i, table := range res.Partition.Tables {
			tables[i] = table.Name
		}

		tw.Append([]string{res.Partition.Name, res.Partition.Description, strings.Join(tables, ",")})
	}

	if tw.NumLines() == 0 {
		return fmt.Errorf("no partitions found in: %s", dir)
	}

	tw.Render()
	return nil
}
