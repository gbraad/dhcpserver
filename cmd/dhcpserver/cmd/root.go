package cmd

import (
	"os"
	"fmt"
	// "github.com/gbraad/dhcpserver/pkg/dhcpserver"
	"github.com/spf13/cobra"
)

const (
	descriptionShort = "Run a simple DHCP server"
	descriptionLong  = "Run a simple DHCP server"
)

var rootCmd = &cobra.Command{
	Use:   commandName,
	Short: descriptionShort,
	Long:  descriptionLong,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		runPrerun()
	},
	Run: func(cmd *cobra.Command, args []string) {
		runRoot()
	},
}

func init() {
	// nothing for now
}

func runPrerun() {
	fmt.Println(commandName)
}

func runRoot() {
	fmt.Println("No command given")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("ERR:", err.Error())
		os.Exit(1)
	}
}
