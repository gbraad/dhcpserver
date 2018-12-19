package cmd

import (
	"github.com/gbraad/dhcpserver/pkg/dhcpserver"
	"github.com/spf13/cobra"
)

var (
	iface string
	port  int
)

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&iface, "interface", "i", "", "Interface to bind to")
	startCmd.Flags().IntVarP(&port, "port", "p", 67, "Port to bind to")
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
	dhcpserver.StartServer(iface, port)
}

