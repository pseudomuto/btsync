package cmd

import (
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

// ConfigFunc represents a function used to configure the CLI
type ConfigFunc func(fs *flag.FlagSet) error

var rootCmd = &cobra.Command{
	Use:     "btsync",
	Short:   "A CLI tool for managing BigTable table definitions",
	Long:    `btsync uses YAML files to describe the desired state of a BigTable cluster's table defintions.`,
	Version: "0.1.0",
}

func init() {
	rootCmd.PersistentFlags().String("dir", "config/partitions", "The directory containing the partition files")
	rootCmd.PersistentFlags().String("instance", "", "The BigTable instance to connect to")
	rootCmd.PersistentFlags().String("project", "", "The project that hosts the BigTable instance")
	rootCmd.PersistentFlags().String("service-account", "", "The service account file to use for auth")
}

// Execute runs the root command object
func Execute(cfn ConfigFunc) error {
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		cfn(rootCmd.PersistentFlags())
	}

	return rootCmd.Execute()
}
