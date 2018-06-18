package cmd

import (
	"encoding/json"
	"fmt"

	"github.build.ge.com/predix-data-services/predix-insights-go-sdk/predixinsights"

	"github.com/spf13/cobra"
)

// flowTemplateCmd represents the flowTemplate command
var flowTemplateCmd = &cobra.Command{
	Use:   "flow-template",
	Short: "Flow Template",
	Long:  `Flow Template.`,
}

func init() {
	RootCmd.AddCommand(flowTemplateCmd)
}

// flowTemplateCmd represents the flowTemplate command
var postFlowTemplateCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a flow template",
	Long:    `Upload a Flow Template to Predix Insights.`,
	Example: `  pi flow-template create --desc "PI CLI Example" --flowTemplateName "pi-cli" --flowType "SPARK_JAVA" --templateFileName "spark-examples.zip" --templateFilePath "/Users/andromeda/Desktop/spark-examples.zip" --flowTemplateVersion 1.0.0`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(postFlowTemplatePI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		ft, err := client.PostFlowTemplate(postFlowTemplatePI.V.GetString("flowTemplateName"), postFlowTemplatePI.V.GetString("templateFileName"), postFlowTemplatePI.V.GetString("templateFilePath"), postFlowTemplatePI.V.GetString("flowTemplateVersion"), postFlowTemplatePI.V.GetString("desc"), postFlowTemplatePI.V.GetString("flowType"))
		if err != nil {
			fmt.Println("error posting flow tempalte err=" + err.Error())
			return
		}
		postFlowTemplatePI.V.Set("flowTemplateID", ft.ID)
		ftByte, _ := json.Marshal(&ft)
		prettyprint(ftByte)
		cleanup(postFlowTemplatePI)
	},
}

var updateFlowTemplateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a flow template",
	Long:    `Update a Predix Insights Flow Template.`,
	Example: `  pi flow-template update --desc "PI CLI Example" --flowTemplateID MY_FLOW_TEMPLATE_ID --flowTemplateName "pi-cli" --flowType "SPARK_JAVA" --templateFileName "spark-examples.zip" --templateFilePath "/Users/andromeda/Desktop/spark-examples.zip" --flowTemplateVersion 1.0.0`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(updateFlowTemplatePI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		err = client.UpdateFlowTemplateByFlowTemplateIDUsingNewZip(updateFlowTemplatePI.V.GetString("flowTemplateID"), updateFlowTemplatePI.V.GetString("flowTemplateName"), updateFlowTemplatePI.V.GetString("templateFileName"), updateFlowTemplatePI.V.GetString("templateFilePath"), updateFlowTemplatePI.V.GetString("flowTemplateVersion"), updateFlowTemplatePI.V.GetString("desc"), updateFlowTemplatePI.V.GetString("flowType"))
		if err != nil {
			fmt.Println("error posting flow tempalte err=" + err.Error())
			return
		}
		fmt.Printf("Successfully updated Flow Template '%s'\n", updateFlowTemplatePI.V.GetString("flowTemplateID"))
		cleanup(updateFlowTemplatePI)
	},
}

var updateFlowTemplateChangeSparkArguments = &cobra.Command{
	Use:     "update-spark-args",
	Short:   "Update Flow Template Spark Arguments",
	Long:    `Update a Predix Insights Flow Template's Spark Arguments.`,
	Example: `  pi flow-template update-spark-args --flowTemplateID MY_FLOW_TEMPLATE_ID --sparkArgs "{\"sparkArguments\": {\"applicationArgs\":[\"100\"],\"className\":\"org.apache.spark.examples.SparkPi\"}}"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(updateFlowTemplateChangeSparkArgumentsPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}

		sparkArgs := &predixinsights.EncapsulatedSparkArgs{}
		err = json.Unmarshal([]byte(updateFlowTemplateChangeSparkArgumentsPI.V.GetString("sparkArgs")), sparkArgs)
		if err != nil {
			fmt.Println("error invalid format for sparkArgs err=" + err.Error())
			return
		}

		err = client.UpdateFlowTemplateByFlowTemplateIDChangeSparkArguments(updateFlowTemplateChangeSparkArgumentsPI.V.GetString("flowTemplateID"), *sparkArgs)
		if err != nil {
			fmt.Println("error updating flow template spark arguments err=" + err.Error())
			return
		}
		fmt.Printf("Successfully updated Flow Template '%s'\n", updateFlowTemplateChangeSparkArgumentsPI.V.GetString("flowTemplateID"))
		cleanup(updateFlowTemplateChangeSparkArgumentsPI)
	},
}

// getFlowTemplateCmd represents the flowTemplate command
var getFlowTemplateCmd = &cobra.Command{
	Use:     "list",
	Short:   "List Flow Template(s)",
	Long:    `List Predix Insights Flow Templates.`,
	Example: "  pi flow-template list --flowTemplateID MY_FLOW_TEMPLATE_ID\n  pi flow-template-list --flowTemplateName MY_FLOW_TEMPLATE_NAME\n  pi flow-template-list",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getFlowTemplatePI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}

		if getFlowTemplatePI.V.GetString("flowTemplateID") != "" {
			flowTemplate, err := client.GetFlowTemplate(getFlowTemplatePI.V.GetString("flowTemplateID"))
			if err != nil {
				fmt.Println("error getting flow tempalte err=" + err.Error())
				return
			}
			ft, _ := json.Marshal(&flowTemplate)
			prettyprint(ft)
		} else if getFlowTemplatePI.V.GetString("flowTemplateName") != "" {
			flowTemplatesResponseWithMetadata, err := client.GetFlowTemplateByName(getFlowTemplatePI.V.GetString("flowTemplateName"))
			if err != nil {
				fmt.Println("error getting flow tempaltes err=" + err.Error())
				return
			}
			ft, _ := json.Marshal(&flowTemplatesResponseWithMetadata)
			prettyprint(ft)

		} else {
			flowTemplatesResponseWithMetadata, err := client.GetAllFlowTemplates()
			if err != nil {
				fmt.Println("error getting all flow tempaltes err=" + err.Error())
				return
			}
			ft, _ := json.Marshal(&flowTemplatesResponseWithMetadata)
			prettyprint(ft)
		}
		cleanup(getFlowTemplatePI)
	},
}

var getFlowTemplateTagsCmd = &cobra.Command{
	Use:     "list-tags",
	Short:   "List Flow Template Tag(s)",
	Long:    `List Predix Insights Flow Template Tag(s).`,
	Example: "  pi flow-template list-tags --flowTemplateID MY_FLOW_TEMPLATE_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getFlowTemplateTagsPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		tagsArray, err := client.GetTagsByFlowTemplateID(saveFlowTemplateTagsPI.V.GetString("flowTemplateID"))
		if err != nil {
			fmt.Println("error getting flow template tags err=" + err.Error())
			return
		}
		tags, _ := json.Marshal(&tagsArray)
		prettyprint(tags)
		cleanup(getFlowTemplateTagsPI)
	},
}

var saveFlowTemplateTagsCmd = &cobra.Command{
	Use:     "save-tags",
	Short:   "Save Flow Template Tag(s)",
	Long:    `Save Predix Insights Flow Template Tag(s).`,
	Example: `  pi flow-template save-tags --flowTemplateID MY_FLOW_TEMPLATE_ID --tags "[\"type:dev\", \"size:large\"]"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(saveFlowTemplateTagsPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		tagsArray := &predixinsights.TagsArray{}
		err = json.Unmarshal([]byte(saveFlowTemplateTagsPI.V.GetString("tags")), tagsArray)
		if err != nil {
			fmt.Println("error invalid format for tags err=" + err.Error())
			return
		}
		saveTagsForFlowTemplateResponse, err := client.SaveTagsForFlowTemplate(saveFlowTemplateTagsPI.V.GetString("flowTemplateID"), *tagsArray)
		if err != nil {
			fmt.Println("error saving flow template tags err=" + err.Error())
			return
		}
		tags, _ := json.Marshal(&saveTagsForFlowTemplateResponse)
		prettyprint(tags)
		cleanup(getFlowTemplateTagsPI)
	},
}

// deleteFlowTemplateCmd represents the flowTemplate command
var deleteFlowTemplateCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a Flow Templates",
	Long:    `Permanently delete a Flow Template.`,
	Example: "  pi flow-template-delete --flowTemplateID MY_FLOW_TEMPLATE_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(deleteFlowTemplatePI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		if !deleteFlowTemplatePI.V.GetBool("force") {
			fmt.Printf("Really delete the Flow Template '%s'? ", deleteFlowTemplatePI.V.GetString("flowTemplateID"))
			if !askForConfirmation() {
				return
			}
		}

		err = client.DeleteFlowTemplate(deleteFlowTemplatePI.V.GetString("flowTemplateID"))
		if err != nil {
			fmt.Println("error posting flow tempalte err=" + err.Error())
			return
		}
		fmt.Printf("Successfully deleted Flow Template '%s'\n", deleteFlowTemplatePI.V.GetString("flowTemplateID"))
		deleteFlowTemplatePI.V.Set("flowTemplateID", "")
		deleteFlowTemplatePI.V.Set("flowTemplateName", "")
		deleteFlowTemplatePI.V.Set("templateFileName", "")
		deleteFlowTemplatePI.V.Set("templateFilePath", "")
		deleteFlowTemplatePI.V.Set("flowTemplateVersion", "")
		deleteFlowTemplatePI.V.Set("desc", "")
		deleteFlowTemplatePI.V.Set("flowType", "")
		cleanup(deleteFlowTemplatePI)
	},
}
