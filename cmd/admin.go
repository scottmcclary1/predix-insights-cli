package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var adminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin",
	Long:  `Admin.`,
}

func init() {
	RootCmd.AddCommand(adminCmd)
}

// healthCheckCmd represents the healthCheck command
var healthCheckCmd = &cobra.Command{
	Use:     "health-check",
	Short:   "Health Check for Predix Insights",
	Long:    `Check the current health of Predix Insights.`,
	Example: "  pi admin health-check",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(healthCheckPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		err = client.CheckStatus()
		if err != nil {
			fmt.Println("health check failed err=" + err.Error())
			return
		}
		fmt.Println("Up and running!")
		cleanup(healthCheckPI)
	},
}

var versionCheckCmd = &cobra.Command{
	Use:     "version",
	Short:   "Predix Insights API Version",
	Long:    `Check the current Predix Insights API Artifact Version.`,
	Example: "  pi admin version",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(healthCheckPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		version, err := client.CheckVersion()
		if err != nil {
			fmt.Println("version check failed err=" + err.Error())
			return
		}
		fmt.Println(version)
		cleanup(versionCheckPI)
	},
}
