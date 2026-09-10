package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/sandeep-yadav/nexus-ast/pkg/web"
)

var port int

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch the interactive web visualizer and authorization graph dashboard",
	Run: func(cmd *cobra.Command, args []string) {
		if err := web.StartServer(port); err != nil {
			log.Fatalf("Failed to launch dashboard: %v", err)
		}
	},
}

func init() {
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to host web dashboard on")
	rootCmd.AddCommand(serveCmd)
}
