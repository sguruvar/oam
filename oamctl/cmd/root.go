/*
Copyright © 2023 Siva Guruvareddiar
*/
package cmd

import (
	"os"

	"github.com/sguruvar/oamctl/cmd/create"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "oamctl",
	Short: "OAM Command Line Tool for Kubevela",
	Long:  `OAM Command Line Tool for Kubevela`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubcommands() {
	rootCmd.AddCommand(create.CreateCmd)

}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addSubcommands()
}
