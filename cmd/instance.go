package cmd

import (
	"encoding/json"
	"fmt"
	"time"

	"github.build.ge.com/predix-data-services/predix-insights-go-sdk/predixinsights"
	tm "github.com/buger/goterm"
	"github.com/spf13/cobra"
)

var instanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "Flow Instance",
	Long:  `Flow Instance.`,
}

func init() {
	RootCmd.AddCommand(instanceCmd)
}

var getInstanceCmd = &cobra.Command{
	Use:     "list-instance",
	Short:   "List Predix Insights Flow Instance(s)",
	Long:    `List Predix Insights Flow Instance(s).`,
	Example: "  pi instance list-instance --instance ID MY_INSTANCE_ID\n  pi instance list-instance",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getInstancePI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		if getInstancePI.V.GetString("instanceID") != "" {
			instanceResponse, err := client.GetInstance(getInstancePI.V.GetString("instanceID"))
			if err != nil {
				fmt.Println("error getting instance err=" + err.Error())
				return
			}
			instanceByte, _ := json.Marshal(&instanceResponse)
			prettyprint(instanceByte)
		} else {

			instanceResponse, err := client.GetAllInstances()
			if err != nil {
				fmt.Println("error getting all instances err=" + err.Error())
				return
			}
			instanceByte, _ := json.Marshal(&instanceResponse)
			prettyprint(instanceByte)
		}
		cleanup(getInstancePI)
	},
}

var getAllInstanceContainers = &cobra.Command{
	Use:     "list-containers",
	Short:   "List Flow Containers",
	Long:    `List Predix Insights Flow Containers.`,
	Example: "  pi instance list-containers --instanceID MY_INSTANCE_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getAllInstanceContainersPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}

		containerResponse, err := client.GetAllInstanceContainers(getAllInstanceContainersPI.V.GetString("instanceID"))
		if err != nil {
			fmt.Println("error getting instance err=" + err.Error())
			return
		}
		containerByte, _ := json.Marshal(&containerResponse)
		prettyprint(containerByte)
		cleanup(getAllInstanceContainersPI)
	},
}

var stopInstanceCmd = &cobra.Command{
	Use:     "stop",
	Short:   "Stop Flow Instance",
	Long:    `Stop Predix Insights Flow Instance.`,
	Example: "  pi instance stop --instanceID MY_INSTANCE_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(stopInstancePI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		err = client.StopInstance(stopInstancePI.V.GetString("instanceID"))
		if err != nil {
			fmt.Println("error stopping instance err=" + err.Error())
			return
		}
		cleanup(stopInstancePI)
	},
}

var getInstanceSubmitLogsCmd = &cobra.Command{
	Use:     "list-submit-logs",
	Short:   "List Flow Instance Submit Logs",
	Long:    `List Predix Insights Flow Instance Submit Logs.`,
	Example: "  pi instance list-submit-logs --instanceID MY_INSTANCE_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getInstanceSubmitLogsPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		submitLogs, err := client.GetInstanceSubmitLogsByInstanceID(getInstanceSubmitLogsPI.V.GetString("instanceID"))
		if err != nil {
			fmt.Println("error getting flow instance submit logs err=" + err.Error())
			return
		}
		fmt.Println(submitLogs)
		cleanup(getInstanceSubmitLogsPI)
	},
}

var getContainerLogsResponse = &cobra.Command{
	Use:     "list-container-response",
	Short:   "List Flow Container Logs Response",
	Long:    `List Predix Insights Flow Container Logs Response.`,
	Example: "  pi instance list-container-response --instanceID MY_INSTANCE_ID --containerID MY_CONTAINER_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getContainerLogsResponsePI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}

		containerLogsResponse, err := client.GetContainerLogsByInstanceIDAndContainerID(getContainerLogsResponsePI.V.GetString("instanceID"), getContainerLogsResponsePI.V.GetString("containerID"))
		if err != nil {
			fmt.Println("error getting container logs response err=" + err.Error())
			return
		}
		containerByte, _ := json.Marshal(&containerLogsResponse)
		prettyprint(containerByte)
		cleanup(getContainerLogsResponsePI)
	},
}

var getContainerLogs = &cobra.Command{
	Use:     "list-container-logs",
	Short:   "List Container Logs",
	Long:    `List Predix Insights Container Logs.`,
	Example: "  pi instance list-container-logs --instanceID MY_INSTANCE_ID --containerID MY_CONTAINER_ID\n  pi instance list-container-logs --instanceID MY_INSTANCE_ID --containerID MY_CONTAINER_ID --tail\n  pi instance list-container-logs --instanceID MY_INSTANCE_ID --containerID MY_CONTAINER_ID --containerLogSink 0\n  pi instance list-container-logs --instanceID MY_INSTANCE_ID --containerID MY_CONTAINER_ID --containerLogSink 1",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getContainerLogsPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		if getContainerLogsPI.V.GetBool("tail") {
			tm.Clear()
			for {
				tm.MoveCursor(1, 1)
				logs, err := client.GetInstanceContainerLogs(getContainerLogsPI.V.GetString("instanceID"), getContainerLogsPI.V.GetString("containerID"), predixinsights.ContainerLogSink(getContainerLogsPI.V.GetInt("containerLogSink")))
				if err != nil {
					fmt.Println("error getting container logs err=" + err.Error())
					return
				}
				tm.Println(logs)
				tm.Flush()
				time.Sleep(time.Second * 5)
			}
		} else {
			logs, err := client.GetInstanceContainerLogs(getContainerLogsPI.V.GetString("instanceID"), getContainerLogsPI.V.GetString("containerID"), predixinsights.ContainerLogSink(getContainerLogsPI.V.GetInt("containerLogSink")))
			if err != nil {
				fmt.Println("error getting container logs err=" + err.Error())
				return
			}
			fmt.Println(logs)
		}
		cleanup(getContainerLogsPI)
	},
}

var getSparkAppDetails = &cobra.Command{
	Use:     "list-spark-app-details",
	Short:   "List Spark Application Details",
	Long:    `List Spark Application Details.`,
	Example: "  pi instance list-spark-app-details --instanceID MY_INSTANCE_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getSparkAppDetailsPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		appDetails, err := client.GetSparkApplicationDetails(getSparkAppDetailsPI.V.GetString("instanceID"))
		if err != nil {
			fmt.Println("error getting spark application details err=" + err.Error())
			return
		}
		appDetailsByte, _ := json.Marshal(&appDetails)
		prettyprint(appDetailsByte)
		cleanup(getSparkAppDetailsPI)

	},
}

var getSparkExecutorDetails = &cobra.Command{
	Use:     "list-spark-executor-details",
	Short:   "List Spark Executor Details",
	Long:    `List Spark Executor Details.`,
	Example: "  pi instance list-spark-executor-details --instanceID MY_INSTANCE_ID --attemptID MY_ATTEMPT_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getSparkExecutorDetailsPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		executorDetails, err := client.GetSparkExecutorDetails(getSparkExecutorDetailsPI.V.GetString("instanceID"), getSparkExecutorDetailsPI.V.GetString("attemptID"))
		if err != nil {
			fmt.Println("error getting spark executor details err=" + err.Error())
			return
		}
		executorDetailsByte, _ := json.Marshal(&executorDetails)
		prettyprint(executorDetailsByte)
		cleanup(getSparkExecutorDetailsPI)
	},
}

var getAllAppStages = &cobra.Command{
	Use:     "list-app-stages",
	Short:   "List Stages of Application Instance",
	Long:    `List Stages of Application Instance.`,
	Example: "  pi instance list-app-stages --instanceID MY_INSTANCE_ID --attemptID MY_ATTEMPT_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getAllAppStagesPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		stageInformation, err := client.GetAllStagesOfApplicationInstance(getAllAppStagesPI.V.GetString("instanceID"), getAllAppStagesPI.V.GetString("attemptID"))
		if err != nil {
			fmt.Println("error getting all stages of application instance err=" + err.Error())
			return
		}
		stageInformationByte, _ := json.Marshal(&stageInformation)
		prettyprint(stageInformationByte)
		cleanup(getAllAppStagesPI)
	},
}

var getAllAttempts = &cobra.Command{
	Use:     "list-attempts",
	Short:   "List All Attempts by Stage",
	Long:    `List All Attempts by Stage.`,
	Example: "  pi instance list-attempts --instanceID MY_INSTANCE_ID --attemptID MY_ATTEMPT_ID --stageID MY_STAGE_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getAllAttemptsPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		allAttemptsForStage, err := client.GetAllAttemptsByStage(getAllAttemptsPI.V.GetString("instanceID"), getAllAttemptsPI.V.GetString("attemptID"), getAllAttemptsPI.V.GetString("stageID"))
		if err != nil {
			fmt.Println("error getting all stages of application instance err=" + err.Error())
			return
		}
		allAttemptsForStageByte, _ := json.Marshal(&allAttemptsForStage)
		prettyprint(allAttemptsForStageByte)
		cleanup(getAllAttemptsPI)
	},
}

var getAttemptDetails = &cobra.Command{
	Use:     "list-attempt-details",
	Short:   "List Stage Attempt Details",
	Long:    `List Stage Attempt Details.`,
	Example: "  pi instance list-attempt-details MY_INSTANCE_ID --attemptID MY_ATTEMPT_ID --stageID MY_STAGE_ID --stageAttemptID MY_STAGE_ATTEMPT_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getAttemptDetailsPI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		allAttemptsForStage, err := client.GetStageAttemptDetails(getAttemptDetailsPI.V.GetString("instanceID"), getAttemptDetailsPI.V.GetString("attemptID"), getAttemptDetailsPI.V.GetString("stageID"), getAttemptDetailsPI.V.GetString("stageAttemptID"))
		if err != nil {
			fmt.Println("error getting stage attempt details err=" + err.Error())
			return
		}
		allAttemptsForStageByte, _ := json.Marshal(&allAttemptsForStage)
		prettyprint(allAttemptsForStageByte)
		cleanup(getAttemptDetailsPI)
	},
}

var getAllTasksByStage = &cobra.Command{
	Use:     "list-tasks",
	Short:   "List All Tasks by Stage",
	Long:    `List All Tasks by Stage.`,
	Example: "  pi instance list-tasks --instanceID MY_INSTANCE_ID --attemptID MY_ATTEMPT_ID --stageID MY_STAGE_ID --stageAttemptID MY_STAGE_ATTEMPT_ID",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		err = getMissingRequiredParams(getAllTasksByStagePI)
		if err != nil {
			fmt.Println("failed to get required parameters err=" + err.Error())
			return
		}
		tasks, err := client.GetAllTasksByStage(getAllTasksByStagePI.V.GetString("instanceID"), getAllTasksByStagePI.V.GetString("attemptID"), getAllTasksByStagePI.V.GetString("stageID"), getAllTasksByStagePI.V.GetString("stageAttemptID"))
		if err != nil {
			fmt.Println("error getting all tasks by stage err=" + err.Error())
			return
		}
		tasksByte, _ := json.Marshal(&tasks)
		prettyprint(tasksByte)
		cleanup(getAllTasksByStagePI)
	},
}
