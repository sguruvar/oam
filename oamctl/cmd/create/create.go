/*
Copyright © 2023 Siva Guruvareddiar
*/
package create

import (
	"github.com/spf13/cobra"
)

// CreateCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create Kubevela deployment YAML ",
	Long:  `create Kubevela deployment YAML `,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {

}
