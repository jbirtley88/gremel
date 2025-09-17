package cmd

import (
	"log"

	"github.com/jbirtley88/gremel/api"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:     "daemon",
	Aliases: []string{"server"},
	Short:   "Runs the Gremel daemon (API server)",
	Run:     RunDaemon,
}

var daemonAddress string

func init() {
	// add the queue subcommand
	rootCmd.AddCommand(daemonCmd)

	daemonCmd.PersistentFlags().StringVarP(&daemonAddress, "listen", "l", "0.0.0.0:8000", "Gremel daemon listen address")
}

func RunDaemon(cmd *cobra.Command, args []string) {
	// Validate args and flags
	if daemonAddress == "" {
		log.Fatal("RunDaemon(): Must have a value for -listen")
	}

	log.Println("Running Gremel daemon on " + daemonAddress)

	// Create the router
	ginRouter, err := api.NewRouter()
	if err != nil {
		log.Fatalf("RunDaemon(): Could not start server: %s", err.Error())
		log.Fatalf("RunDaemon(): Could not start server: %s", err.Error())
	}

	// Set the route handlers
	ginRouter.PUT("/api/v1/mount", api.MountTable)
	ginRouter.GET("/api/v1/mount", api.GetMount)
	ginRouter.GET("/api/v1/query", api.Query)
	ginRouter.GET("/api/v1/schema", api.Schema)
	ginRouter.GET("/api/v1/tables", api.Tables)

	// Start the service
	err = ginRouter.Run(daemonAddress)
	if err != nil {
		log.Fatalf("RunDaemon(): Could not start server: %s", err.Error())
	}
}
