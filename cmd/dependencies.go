package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(dependencyCmd)
}

var dependencyCmd = &cobra.Command{
	Use:   "dependency",
	Short: "Dependency",
	Long:  `Predix Insights Dependency.`,
}

var getDependencyCmd = &cobra.Command{
	Use:     "list",
	Short:   "List Dependencies",
	Long:    `List Predix Insights Dependencies.`,
	Example: "  pi dependency list --dependencyID MY_DEPENDENCY_ID\n  pi dependency list",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getDependencyPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		if getDependencyPI.V.GetString("dependencyID") != "" {
			flow, err := client.GetDependencyByID(getDependencyPI.V.GetString("dependencyID"))
			if err != nil {
				fmt.Println("error getting dependency err=" + err.Error())
				return
			}
			flowByte, _ := json.Marshal(&flow)
			prettyprint(flowByte)
		} else {
			flows, err := client.GetAllDependencies()
			if err != nil {
				fmt.Println("error getting all dependencies err=" + err.Error())
				return
			}
			flowsByte, _ := json.Marshal(&flows)
			prettyprint(flowsByte)
		}
		cleanup(getDependencyPI)
	},
}

var deployDependencyCmd = &cobra.Command{
	Use:     "deploy",
	Short:   "Deploy Dependencies",
	Long:    `Deploy Predix Insights Dependencies.`,
	Example: "  pi dependency deploy --dependencyID MY_DEPENDENCY_ID\n  pi dependency deploy",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(deployDependencyPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		if deployDependencyPI.V.GetString("dependencyID") != "" {
			err = client.DeployDependencyByDependencyID(deployDependencyPI.V.GetString("dependencyID"))
			if err != nil {
				fmt.Println("error deploying dependency err=" + err.Error())
				return
			}
		} else {
			err = client.DeployAllDependencies()
			if err != nil {
				fmt.Println("error deploying all dependencies err=" + err.Error())
				return
			}
		}
		cleanup(deployDependencyPI)
	},
}

var unDeployDependencyCmd = &cobra.Command{
	Use:     "undeploy",
	Short:   "Undeploy Dependencies",
	Long:    `Undeploy Predix Insights Dependencies.`,
	Example: "  pi dependency deploy --dependencyID MY_DEPENDENCY_ID\n  pi dependency deploy",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(unDeployDependencyPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		if unDeployDependencyPI.V.GetString("dependencyID") != "" {
			err = client.UnDeployDependencyByDependencyID(unDeployDependencyPI.V.GetString("dependencyID"))
			if err != nil {
				fmt.Println("error undeploying dependency err=" + err.Error())
				return
			}
		} else {
			err = client.UnDeployAllDependencies()
			if err != nil {
				fmt.Println("error undeploying all dependencies err=" + err.Error())
				return
			}
		}
		cleanup(unDeployDependencyPI)
	},
}

var deleteDependencyCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a Dependency",
	Long:    `Permanently delete a Dependency.`,
	Example: "  pi dependency delete --dependencyID MY_DEPENDENCY_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(deleteDependencyPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		if !deleteDependencyPI.V.GetBool("force") {
			fmt.Printf("Really delete the Dependency '%s'? ", deleteDependencyPI.V.GetString("dependencyID"))
			if !askForConfirmation() {
				return
			}
		}
		err = client.DeleteDependencyByID(deleteDependencyPI.V.GetString("dependencyID"))
		if err != nil {
			fmt.Println("error deleting dependency err=" + err.Error())
			return
		}
		deleteDependencyPI.V.Set("dependencyID", "")
		cleanup(deleteDependencyPI)
	},
}

var postDependencyCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a Dependency",
	Long:    `Create a Predix Insights Dependency.`,
	Example: "  pi dependency create --dependencyType MY_DEPENDENCY_TYPE --dependencyFileName MY_DEPENDENCY_FILE_NAME --dependencyFileLocation MY_DEPENDENCY_FILE_LOCATION",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(postDependencyPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		dependencyResponse, err := client.PostDependency(postDependencyPI.V.GetString("dependencyType"), postDependencyPI.V.GetString("dependencyFileName"), postDependencyPI.V.GetString("dependencyFileLocation"))
		if err != nil {
			fmt.Println("error deleting dependency err=" + err.Error())
			return
		}
		dependencyResponseByte, _ := json.Marshal(&dependencyResponse)
		prettyprint(dependencyResponseByte)
		if len(dependencyResponse) > 0 {
			postDependencyPI.V.Set("dependencyID", dependencyResponse[0].ID)
			postDependencyPI.V.Set("dependencyName", dependencyResponse[0].Name)
		}
		cleanup(postDependencyPI)
	},
}
