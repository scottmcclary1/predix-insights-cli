package cmd

import (
	"encoding/json"
	"fmt"

	"github.build.ge.com/predix-data-services/predix-insights-go-sdk/predixinsights"
	"github.com/spf13/cobra"
)

// flowCmd represents the flow command
var flowCmd = &cobra.Command{
	Use:   "flow",
	Short: "Flow",
	Long:  `Flow.`,
}

func init() {
	RootCmd.AddCommand(flowCmd)
}

// flowCmd represents the flow command
var getFlowCmd = &cobra.Command{
	Use:     "list",
	Short:   "List Flow(s)",
	Long:    `List Predix Insights Flow(s).`,
	Example: "  pi flow list --flowName MY_FLOW_NAME\n pi flow list --flowID MY_FLOW_ID\n pi flow list --flowTemplateID MY_FLOW_TEMPLATE_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getFlowPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		if getFlowPI.V.GetString("flowName") != "" {
			flow, err := client.GetFlow(getFlowPI.V.GetString("flowName"))
			if err != nil {
				fmt.Println("error getting all flow err=" + err.Error())
				return
			}
			flowByte, _ := json.Marshal(&flow)
			prettyprint(flowByte)
		} else if getFlowPI.V.GetString("flowID") != "" {
			flowResponse, err := client.GetFlowByTemplateIDAndFlowID(getFlowPI.V.GetString("flowTemplateID"), getFlowPI.V.GetString("flowID"))
			if err != nil {
				fmt.Println("error getting flow err=" + err.Error())
				return
			}
			flow, _ := json.Marshal(&flowResponse)
			prettyprint(flow)
		} else if getFlowPI.V.GetString("flowTemplateID") != "" {
			getAllFlowsByTemplateIDResponse, err := client.GetAllFlowsByTemplateID(getFlowPI.V.GetString("flowTemplateID"))
			if err != nil {
				fmt.Println("error getting flows err=" + err.Error())
				return
			}
			flows, _ := json.Marshal(&getAllFlowsByTemplateIDResponse)
			prettyprint(flows)
		} else {
			flows, err := client.GetAllFlows(1)
			if err != nil {
				fmt.Println("error getting all flows err=" + err.Error())
				return
			}
			flowsByte, _ := json.Marshal(&flows)
			prettyprint(flowsByte)
		}
		cleanup(getFlowPI)
	},
}

// deleteFlowCmd represents the flow command
var deleteFlowCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a Flow",
	Long:    `Permanently delete a Flow.`,
	Example: "  pi flow delete --flowID MY_FLOW_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(deleteFlowPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		if !deleteFlowPI.V.GetBool("force") {
			fmt.Printf("Really delete the flow '%s'? ", deleteFlowPI.V.GetString("flowID"))
			if !askForConfirmation() {
				return
			}
		}
		err = client.DeleteFlowByFlowIDOnly(deleteFlowPI.V.GetString("flowID"))
		if err != nil {
			fmt.Println("error deleting flow err=" + err.Error())
			return
		}
		fmt.Printf("Successfully deleted Flow '%s'\n", deleteFlowPI.V.GetString("flowID"))
		deleteFlowPI.V.Set("flowID", "")
		deleteFlowPI.V.Set("flowName", "")
		cleanup(deleteFlowPI)
	},
}

// postFlowCmd represents the flow command
var postFlowCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a Flow",
	Long:    `Create a Predix Insights Flow.`,
	Example: "  pi flow delete --flowName MY_FLOW_NAME --flowTemplateID MY_FLOW_TEMPLATE_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(postFlowPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		flow, err := client.PostFlow(postFlowPI.V.GetString("flowName"), postFlowPI.V.GetString("flowTemplateID"))
		if err != nil {
			fmt.Println("error posting flow err=" + err.Error())
			return
		}

		postFlowPI.V.Set("flowID", flow.ID)
		postFlowPI.V.Set("flowName", flow.Name)
		postFlowByte, _ := json.Marshal(&flow)
		prettyprint(postFlowByte)
		cleanup(postFlowPI)
	},
}

var postDirectFlowCmd = &cobra.Command{
	Use:     "create-direct",
	Short:   "Create a Direct Flow",
	Long:    `Create a Predix Insights Direct Flow.`,
	Example: "  pi flow create-direct --flowName MY_FLOW_NAME --flowFileName test.zip --flowFilePath /Users/andromeda/Desktop/test.zip --flowVersion 1.0.0 --desc \"My description\" --flowType SPARK_JAVA",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(postDirectFlowPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		flow, err := client.PostFlowDirectly(postDirectFlowPI.V.GetString("flowName"), postDirectFlowPI.V.GetString("flowFileName"), postDirectFlowPI.V.GetString("flowFilePath"), postDirectFlowPI.V.GetString("flowVersion"), postDirectFlowPI.V.GetString("desc"), postDirectFlowPI.V.GetString("flowType"))
		if err != nil {
			fmt.Println("error posting direct flow err=" + err.Error())
			return
		}

		postDirectFlowPI.V.Set("flowID", flow.ID)
		postDirectFlowPI.V.Set("flowName", flow.Name)
		postFlowByte, _ := json.Marshal(&flow)
		prettyprint(postFlowByte)
		cleanup(postDirectFlowPI)
	},
}

var updateDirectFlowCmd = &cobra.Command{
	Use:     "update-direct",
	Short:   "Update a Direct Flow",
	Long:    `Update a Predix Insights Direct Flow.`,
	Example: "  pi flow update-direct --flowID MY_FLOW_ID --flowFileName test.zip --flowFilePath /Users/andromeda/Desktop/test.zip --desc \"My description\"",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(updateDirectFlowPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		flow, err := client.UpdateDirectFlowByFlowIDChangeAnalyticFile(updateDirectFlowPI.V.GetString("flowID"), updateDirectFlowPI.V.GetString("desc"), updateDirectFlowPI.V.GetString("flowFileName"), updateDirectFlowPI.V.GetString("flowFilePath"))
		if err != nil {
			fmt.Println("error updating direct flow err=" + err.Error())
			return
		}

		updateDirectFlowPI.V.Set("flowID", flow.ID)
		updateDirectFlowPI.V.Set("flowName", flow.Name)
		postFlowByte, _ := json.Marshal(&flow)
		prettyprint(postFlowByte)
		cleanup(updateDirectFlowPI)
	},
}

var postLaunchFlowCmd = &cobra.Command{
	Use:     "launch",
	Short:   "Launch a Flow",
	Long:    `Launch a Predix Insights Flow.`,
	Example: "  pi flow launch --flowID MY_FLOW_ID --flowTemplateID MY_FLOW_TEMPLATE_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(postLaunchFlowPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		launchResponse, err := client.LaunchFlow(postLaunchFlowPI.V.GetString("flowTemplateID"), postLaunchFlowPI.V.GetString("flowID"))
		if err != nil {
			fmt.Println("error launching flow err=" + err.Error())
			return
		}
		postLaunchFlowPI.V.Set("instanceID", launchResponse.ID)
		postFlowByte, _ := json.Marshal(&launchResponse)
		prettyprint(postFlowByte)
		cleanup(postLaunchFlowPI)
	},
}

var stopFlowCmd = &cobra.Command{
	Use:     "stop",
	Short:   "Stop a Flow",
	Long:    `Stop a Predix Insights Flow.`,
	Example: "  pi flow stop --flowName MY_FLOW_NAME",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(stopFlowPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		err = client.StopFlow(stopFlowPI.V.GetString("flowName"))
		if err != nil {
			fmt.Println("error stopping flow err=" + err.Error())
			return
		}
		fmt.Printf("Flow %s successfully stoppped.\n", stopFlowPI.V.GetString("flowName"))
		cleanup(stopFlowPI)
	},
}

var createFlowTemplateFromFlowCmd = &cobra.Command{
	Use:     "create-flow-template",
	Short:   "Create a Flow Template",
	Long:    `Create a Predix Insights Flow Template from a Direct Flow.`,
	Example: "  pi flow create-flow-template --flowID MY_FLOW_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(createFlowTemplateFromFlowPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		ft, err := client.CreateFlowTemplateFromFlow(createFlowTemplateFromFlowPI.V.GetString("flowID"))
		if err != nil {
			fmt.Println("error stopping flow err=" + err.Error())
			return
		}
		createFlowTemplateFromFlowPI.V.Set("flowTemplateID", ft.ID)
		createFlowTemplateFromFlowPI.V.Set("flowTemplateName", ft.Name)
		ftByte, _ := json.Marshal(&ft)
		prettyprint(ftByte)
		cleanup(createFlowTemplateFromFlowPI)
	},
}

var updateFlowChangeSparkArguments = &cobra.Command{
	Use:     "update-spark-args",
	Short:   "Update Flow Spark Arguments",
	Long:    `Update a Predix Insights Flow's Spark Arguments.`,
	Example: `  pi flow update-spark-args --sparkArgs "{\"sparkArguments\": {\"applicationArgs\":[\"100\"],\"className\":\"org.apache.spark.examples.SparkPi\"}}" -i\n pi flow update-spark-args --sparkArgs "{\"sparkArguments\": {\"applicationArgs\":[\"1000\"],\"className\":\"org.apache.spark.examples.SparkPi\"}}" --flowTemplateID MY_FLOW_TEMPLATE_ID --flowID MY_FLOW_ID`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(updateFlowChangeSparkArgumentsPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		sparkArgs := &predixinsights.EncapsulatedSparkArgs{}
		err = json.Unmarshal([]byte(updateFlowChangeSparkArgumentsPI.V.GetString("sparkArgs")), sparkArgs)
		if err != nil {
			fmt.Println("error invalid format for sparkArgs err=" + err.Error())
			return
		}

		err = client.UpdateFlowChangeSparkArguments(updateFlowChangeSparkArgumentsPI.V.GetString("flowTemplateID"), updateFlowChangeSparkArgumentsPI.V.GetString("flowID"), *sparkArgs)
		if err != nil {
			fmt.Println("error updating flow spark arguments err=" + err.Error())
			return
		}
		cleanup(updateFlowChangeSparkArgumentsPI)
		fmt.Printf("Successfully updated Flow '%s'\n", updateFlowChangeSparkArgumentsPI.V.GetString("flowID"))
	},
}

var addFlowConfigFiles = &cobra.Command{
	Use:     "add-config-file",
	Short:   "Add Config File(s) to a Flow",
	Long:    `Add configuration file(s) to a Predix Insights Flow.`,
	Example: `  pi flow add-config-file --flowID MY_FLOW_ID --configFileDetails "[{\"FileName\": \"config.json\", \"FileLocation\": \"/Users/andromeda/Desktop/config.json\"}]"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(addFlowConfigFilesPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		fileDetails := []predixinsights.FileDetails{}
		err = json.Unmarshal([]byte(addFlowConfigFilesPI.V.GetString("configFileDetails")), &fileDetails)
		if err != nil {
			fmt.Println("failed to parse configFileDetails err=" + err.Error())
			return
		}

		err = client.UpdateFlowByFlowIDAddConfigFile(addFlowConfigFilesPI.V.GetString("flowID"), fileDetails)
		if err != nil {
			fmt.Println("error adding config file(s) to flow err=" + err.Error())
			return
		}
		fmt.Printf("Config file(s) successfully added to flow %s.\n", addFlowConfigFilesPI.V.GetString("flowID"))
		cleanup(addFlowConfigFilesPI)
	},
}

var deleteFlowConfigFile = &cobra.Command{
	Use:     "delete-config-file",
	Short:   "Delete a Config File From a Flow",
	Long:    `Delete a configuration file from a Predix Insights Flow.`,
	Example: "  pi flow delete-config-file --flowID MY_FLOW_ID --configFileName MY_CONFIG_FILE_NAME",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(deleteFlowConfigFilePI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		err = client.UpdateFlowByFlowIDDeleteConfigFile(deleteFlowConfigFilePI.V.GetString("flowID"), deleteFlowConfigFilePI.V.GetString("configFileName"))
		if err != nil {
			fmt.Println("error deleting config file(s) err=" + err.Error())
			return
		}
		fmt.Printf("Config file %s successfully deleted from flow %s.\n", deleteFlowConfigFilePI.V.GetString("configFileName"), deleteFlowConfigFilePI.V.GetString("flowID"))
		cleanup(deleteFlowConfigFilePI)
	},
}

var listConfigFilesCmd = &cobra.Command{
	Use:     "list-config-files",
	Short:   "List Predix Insights Flow Config File(s)",
	Long:    `List Predix Insights Flow Configuration File(s).`,
	Example: "  pi flow list-config-files --flowID MY_FLOW_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(listConfigFilesPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		listConfigFiles, err := client.ListConfigFilesByFlowID(listConfigFilesPI.V.GetString("flowID"))
		if err != nil {
			fmt.Println("error getting flow configuration files err=" + err.Error())
			return
		}
		configFilesByte, _ := json.Marshal(&listConfigFiles)
		prettyprint(configFilesByte)
		cleanup(listConfigFilesPI)
	},
}

var saveFlowTagsCmd = &cobra.Command{
	Use:     "save-tags",
	Short:   "Save Flow Tag(s)",
	Long:    `Save Predix Insights Flow Tag(s).`,
	Example: `  pi flow save-tags --flowTemplateID MY_FLOW_TEMPLATE_ID --flowID MY_FLOW_ID --tags "[\"type:dev\", \"size:large\"]"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(saveFlowTagsPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		tagsArray := &predixinsights.TagsArray{}
		err = json.Unmarshal([]byte(saveFlowTagsPI.V.GetString("tags")), tagsArray)
		if err != nil {
			fmt.Println("error invalid format for tags err=" + err.Error())
			return
		}
		flowResponse, err := client.SaveTagsForFlow(saveFlowTagsPI.V.GetString("flowTemplateID"), saveFlowTagsPI.V.GetString("flowID"), *tagsArray)
		if err != nil {
			fmt.Println("error saving flow tags err=" + err.Error())
			return
		}
		tags, _ := json.Marshal(&flowResponse)
		prettyprint(tags)
		cleanup(saveFlowTagsPI)
	},
}

var getFlowTagsCmd = &cobra.Command{
	Use:     "list-tags",
	Short:   "List Flow Tag(s)",
	Long:    `List Predix Insights Flow Tag(s).`,
	Example: "  pi flow list-tags --flowTemplateID MY_FLOW_TEMPLATE_ID --flowID MY_FLOW_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getFlowTagsPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		tagsArray, err := client.GetTagsForFlowByFlowTemplateIDAndFlowID(getFlowTagsPI.V.GetString("flowTemplateID"), getFlowTagsPI.V.GetString("flowID"))
		if err != nil {
			fmt.Println("error saving flow template tags err=" + err.Error())
			return
		}
		tags, _ := json.Marshal(&tagsArray)
		prettyprint(tags)
		cleanup(getFlowTagsPI)
	},
}
