package cmd

import (
	"github.com/gbraad/dhcpserver/pkg/dhcpserver"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start server",
	Long:  "Start server",
	Run: func(cmd *cobra.Command, args []string) {
		runStart(args)
	},
}

func runStart(arguments []string) {
	dhcpserver.StartServer()
}

