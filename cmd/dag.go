package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.build.ge.com/predix-data-services/predix-insights-go-sdk/predixinsights"

	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(dagCmd)
}

var dagCmd = &cobra.Command{
	Use:   "dag",
	Short: "Directed Acyclic Graph (DAG)",
	Long:  `Directed Acyclic Graph (DAG.`,
}

var getDagCmd = &cobra.Command{
	Use:     "list",
	Short:   "List DAG(s)",
	Long:    `List Predix Insights DAG(s).`,
	Example: "  pi dag list --dagName MY_DAG_NAME\n  pi dag list",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getDagPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		if getDagPI.V.GetString("dagName") != "" {
			dag, err := client.GetDAG(getDagPI.V.GetString("dagName"))
			if err != nil {
				fmt.Println("error getting dag err=" + err.Error())
				return
			}
			dagByte, _ := json.Marshal(&dag)
			prettyprint(dagByte)
		} else {
			dags, err := client.GetAllDAGs()
			if err != nil {
				fmt.Println("error getting all dags err=" + err.Error())
				return
			}
			dagsByte, _ := json.Marshal(&dags)
			prettyprint(dagsByte)
		}
		cleanup(getDagPI)
	},
}

var deleteDagCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Delete a DAG",
	Long:    `Permanently delete a DAG.`,
	Example: "  pi dag delete --dagName MY_DAG_NAME",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(deleteDagPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		if !deleteDagPI.V.GetBool("force") {
			fmt.Printf("Really delete the DAG '%s'? ", deleteDagPI.V.GetString("dagName"))
			if !askForConfirmation() {
				return
			}
		}
		err = client.DeleteDAG(deleteDagPI.V.GetString("dagName"))
		if err != nil {
			fmt.Println("error deleting dag err=" + err.Error())
			return
		}
		deleteDagPI.V.Set("dagName", "")
		deleteDagPI.V.Set("dagID", "")
		cleanup(deleteDagPI)
	},
}

var postDagCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a DAG",
	Long:    `Create a Predix Insights DAG.`,
	Example: `pi dag create --dagName MY_DAG_NAME --dagFileName dag.py --dagFilePath /Users/andromeda/Desktop/dag.py --dagVersion 1.0.0 --dagDesc "My dag description" --dagFlowType SPARK_JAVA --dagTemplate "{\"Owner\": \"EMR-SA-86d0-4c60-a462-70f27b9d01c6\", \"FlowName\": \"MyFlowName\", \"Interval\": \"5\"}"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(postDagPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		dt := &predixinsights.DAGTemplate{}
		err = json.Unmarshal([]byte(postDagPI.V.GetString("dagTemplate")), dt)
		dag, err := client.PostDAG(postDagPI.V.GetString("dagName"), postDagPI.V.GetString("dagFileName"), postDagPI.V.GetString("dagFilePath"), postDagPI.V.GetString("dagVersion"), postDagPI.V.GetString("dagDesc"), postDagPI.V.GetString("dagFlowType"), *dt)
		if err != nil {
			fmt.Println("error posting dag err=" + err.Error())
			return
		}

		postDagPI.V.Set("dagID", dag.ID)
		postDagPI.V.Set("dagName", dag.Name)
		postDagByte, _ := json.Marshal(&dag)
		prettyprint(postDagByte)
		cleanup(postDagPI)
	},
}

var updateDagCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update a DAG",
	Long:    `Update a Predix Insights DAG.`,
	Example: `pi dag update --dagName MY_DAG_NAME --dagFileName dag.py --dagFilePath /Users/andromeda/Desktop/dag.py --dagVersion 1.0.0 --dagDesc "My dag description" --dagFlowType SPARK_JAVA --dagTemplate "{\"Owner\": \"EMR-SA-86d0-4c60-a462-70f27b9d01c6\", \"FlowName\": \"MyFlowName\", \"Interval\": \"5\"}"`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(updateDagPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		dt := &predixinsights.DAGTemplate{}
		err = json.Unmarshal([]byte(viper.GetString("dagTemplate")), dt)
		if err != nil {
			fmt.Println("failed to parse dagTemplate err=" + err.Error())
			return
		}
		err = client.UpdateDAG(viper.GetString("dagName"), viper.GetString("dagFileName"), viper.GetString("dagFilePath"), viper.GetString("dagVersion"), viper.GetString("dagDesc"), viper.GetString("dagFlowType"), *dt)
		if err != nil {
			fmt.Println("error updating dag err=" + err.Error())
			return
		}

		fmt.Printf("DAG %s updated successfully\n", viper.GetString("dagName"))
		cleanup(updateDagPI)
	},
}

var deployDagCmd = &cobra.Command{
	Use:     "deploy",
	Short:   "Deploy a DAG",
	Long:    `Deploy a Predix Insights DAG.`,
	Example: "  pi dag deploy --dagName MY_DAG_NAME\n  pi dag deploy",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(deployDagPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		err = client.DeployDAG(deployDagPI.V.GetString("dagName"))
		if err != nil {
			fmt.Println("error deploying dag err=" + err.Error())
			return
		}
		fmt.Printf("DAG %s deployed successfully\n", viper.GetString("dagName"))
		cleanup(deployDagPI)
	},
}

var dagStatusCmd = &cobra.Command{
	Use:     "status",
	Short:   "List DAG(s) Status",
	Long:    `List Predix Insights DAG(s) Status.`,
	Example: "  pi dag status --dagName MY_DAG_NAME\n  pi dag status",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(dagStatusPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		if dagStatusPI.V.GetString("dagName") != "" {
			dag, err := client.GetDAGStatusByDAGName(dagStatusPI.V.GetString("dagName"))
			if err != nil {
				fmt.Println("error getting dag err=" + err.Error())
				return
			}
			dagByte, _ := json.Marshal(&dag)
			prettyprint(dagByte)
		} else {
			dags, err := client.GetAllDAGsAllStatuses()
			if err != nil {
				fmt.Println("error getting all dags err=" + err.Error())
				return
			}
			dagsByte, _ := json.Marshal(&dags)
			prettyprint(dagsByte)
		}
		cleanup(dagStatusPI)
	},
}

var getDagRunCmd = &cobra.Command{
	Use:     "list-run",
	Short:   "List DAG Run(s)",
	Long:    `List Predix Insights DAG Run(s).`,
	Example: "  pi dag list-run --dagName MY_DAG_NAME --dagRunID MY_DAG_RUN_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getDagRunPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		if getDagRunPI.V.GetString("dagRunID") != "" {
			dagRun, err := client.GetRunByDAGNameAndRunID(getDagRunPI.V.GetString("dagName"), getDagRunPI.V.GetString("dagRunID"))
			if err != nil {
				fmt.Println("error getting dag run err=" + err.Error())
				return
			}
			dagRunByte, _ := json.Marshal(&dagRun)
			prettyprint(dagRunByte)
		} else {
			dagRuns, err := client.GetRunsByDAGName(getDagRunPI.V.GetString("dagName"))
			if err != nil {
				fmt.Println("error getting dag runs err=" + err.Error())
				return
			}
			dagRunsByte, _ := json.Marshal(&dagRuns)
			prettyprint(dagRunsByte)
		}
		cleanup(getDagRunPI)
	},
}

var getDagTaskCmd = &cobra.Command{
	Use:     "list-task",
	Short:   "List DAG Task(s)",
	Long:    `List Predix Insights DAG Task(s).`,
	Example: "  pi dag list-task --dagName MY_DAG_NAME --dagTaskID MY_DAG_TASK_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getDagTaskPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		if getDagTaskPI.V.GetString("dagTaskID") != "" {
			dagTask, err := client.GetAllTasksByDagNameAndTaskID(getDagTaskPI.V.GetString("dagName"), getDagTaskPI.V.GetString("dagTaskID"))
			if err != nil {
				fmt.Println("error getting dag task err=" + err.Error())
				return
			}
			dagTaskByte, _ := json.Marshal(&dagTask)
			prettyprint(dagTaskByte)
		} else {
			dagTasks, err := client.GetAllTasksByDagName(getDagTaskPI.V.GetString("dagName"))
			if err != nil {
				fmt.Println("error getting dag tasks err=" + err.Error())
				return
			}
			dagTasksByte, _ := json.Marshal(&dagTasks)
			prettyprint(dagTasksByte)
		}
		cleanup(getDagTaskPI)
	},
}

var getDagTaskRunCmd = &cobra.Command{
	Use:     "task-run-info",
	Short:   "List DAG Task Run Info",
	Long:    `List Predix Insights DAG Task Run Information.`,
	Example: "  pi dag task-run-info --dagName MY_DAG_NAME --dagTaskID MY_DAG_TASK_ID --dagRunID MY_DAG_RUN_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getDagTaskRunPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		dagTaskRun, err := client.GetTaskRunInfo(getDagTaskRunPI.V.GetString("dagName"), getDagTaskRunPI.V.GetString("dagTaskID"), getDagTaskRunPI.V.GetString("dagRunID"))
		if err != nil {
			fmt.Println("error getting dag task run info err=" + err.Error())
			return
		}
		dagTaskRunByte, _ := json.Marshal(&dagTaskRun)
		prettyprint(dagTaskRunByte)
		cleanup(getDagTaskRunPI)
	},
}
