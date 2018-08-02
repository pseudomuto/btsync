package main

import (
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Configure parses the args and configures the global viper instance
func Configure(fs *flag.FlagSet) error {
	viper.AddConfigPath(".")
	viper.SetConfigName("btsync")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("btsync")

	viper.BindPFlags(fs)

	return viper.ReadInConfig()
}
