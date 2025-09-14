package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gremel",
	Short: "Gremel is a swiss-army knife for interrogating disparate data sources with SQL",
	Run: func(cmd *cobra.Command, args []string) {
		RunSQL(cmd, args)
	},
}

func init() {
	cobra.OnInitialize(initialise)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

var cfgFile string
var configOverrides []string

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "config.yml", "config file (default is config.yml)")
	rootCmd.PersistentFlags().StringArrayVarP(&configOverrides, "set", "s", []string{}, "Override config values (e.g. 'config.loglevel=debug')")
}
