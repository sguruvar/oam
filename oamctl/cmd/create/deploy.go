/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Generate OAM Deployment",
	Long:  `Generate OAM deployment file with the values provided`,
	Run: func(cmd *cobra.Command, args []string) {

		s1 := OamDeploy{
			APIVersion: "core.oam.dev/v1beta1",
			Kind:       "Application",
			Metadata: Metadata{
				Name: appName,
			},
			Spec: Spec{
				Components: []Component{{
					Name: appName,
					Type: "webservice",
					Properties: ComponentProperties{
						Image: imagePath,
						Ports: []Port{{
							Port:   port,
							Expose: true,
						},
						},
					},
					Traits: []Trait{{
						Type: "scaler",
						Properties: TraitProperties{
							Replicas: replicas,
						},
					},
					},
				},
				},
			},
		}

		yamlData, err := yaml.Marshal(&s1)

		if err != nil {
			fmt.Printf("Error while Marshaling. %v", err)
		}
		fmt.Println("------------------------------------------------------")
		fmt.Println("Your YAML as follows (available as OAMDeployment.yaml)")
		fmt.Println("------------------------------------------------------")
		fmt.Println(string(yamlData))
		fileName := "OAMDeployment.yaml"
		err = ioutil.WriteFile(fileName, yamlData, 0644)
		if err != nil {
			panic("Unable to write data into the file")
		}

	},
}
var (
	imagePath string
	appName   string
	port      int64
	replicas  int64
)

func init() {

	deployCmd.Flags().StringVarP(&appName, "appName", "a", "", "Application name")
	deployCmd.Flags().Int64VarP(&port, "port", "p", 80, "Port")
	deployCmd.Flags().Int64VarP(&replicas, "replicas", "r", 1, "Number of replicas")
	deployCmd.Flags().StringVarP(&imagePath, "imagePath", "i", "", "Image to pull")

	if err := deployCmd.MarkFlagRequired("imagePath"); err != nil {
		fmt.Println(err)
	}
	if err := deployCmd.MarkFlagRequired("appName"); err != nil {
		fmt.Println(err)
	}
	if err := deployCmd.MarkFlagRequired("port"); err != nil {
		fmt.Println(err)
	}
	if err := deployCmd.MarkFlagRequired("replicas"); err != nil {
		fmt.Println(err)
	}

	CreateCmd.AddCommand(deployCmd)

}
