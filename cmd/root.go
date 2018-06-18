package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "pi",
	Short: "Predix Insights CLI",
	Long: `PI is a CLI library for Predix Insights. This top technology lets you 
concentrate on building analytical pipelines, not managing infrastructure.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Initialize Configuration File
	cobra.OnInitialize(initConfig)

	// Configure Commands
	loginPI = NewPI(
		RootCmd,
		loginCmd,
		[]stringVar{
			stringVar{&APIHost, "APIHost", "", "", "Predix Insights API Host", "API_HOST", true},
			stringVar{&TenantID, "TenantID", "", "", "Predix Insights Tenant ID", "TENANT_ID", true},
			stringVar{&IssuerID, "IssuerID", "", "", "UAA Issuer ID", "ISSUER_ID", true},
			stringVar{&ClientID, "ClientID", "", "", "UUA Client ID", "CLIENT_ID", true},
			stringVar{&ClientSecret, "ClientSecret", "", "", "UAA Client Secret", "CLIENT_SECRET", true},
			stringVar{&Token, "Token", "", "", "UAA Authentication Token", "TOKEN", false},
		},
		[]boolVar{},
		[]intVar{})

	// ADMIN Commands
	// health-check
	healthCheckPI = NewPI(adminCmd, healthCheckCmd, []stringVar{}, []boolVar{}, []intVar{})

	// version
	versionCheckPI = NewPI(adminCmd, versionCheckCmd, []stringVar{}, []boolVar{}, []intVar{})

	// DAG Commands
	// list
	getDagPI = NewPI(
		dagCmd,
		getDagCmd,
		[]stringVar{
			stringVar{&dagName, "dagName", "", "", "DAG Name", "DAG_NAME", false},
		},
		[]boolVar{},
		[]intVar{})
	// delete
	deleteDagPI = NewPI(
		dagCmd,
		deleteDagCmd,
		[]stringVar{stringVar{&dagName, "dagName", "", "", "DAG Name", "DAG_NAME", true}},
		[]boolVar{boolVar{&force, "force", "f", false, "Permanently remove a DAG", "FORCE", false}},
		[]intVar{})
	// create
	postDagPI = NewPI(
		dagCmd,
		postDagCmd,
		[]stringVar{
			stringVar{&dagName, "dagName", "", "", "DAG Name", "DAG_NAME", true},
			stringVar{&dagFileName, "dagFileName", "", "", "DAG File Name", "DAG_FILE_NAME", true},
			stringVar{&dagFilePath, "dagFilePath", "", "", "DAG File Path", "DAG_FILE_PATH", true},
			stringVar{&dagVersion, "dagVersion", "", "", "DAG Version", "DAG_VERSION", true},
			stringVar{&dagDesc, "dagDesc", "", "", "DAG Description", "DAG_DESCRIPTION", true},
			stringVar{&dagFlowType, "dagFlowType", "", "", "DAG Flow Type (SPARK_JAVA or SPARK_PYTHON)", "DAG_FLOW_TYPE", true},
			stringVar{&dagTemplate, "dagTemplate", "", "", "DAG Template", "DAG_TEMPLATE", true},
		},
		[]boolVar{},
		[]intVar{})
	// update
	updateDagPI = NewPI(
		dagCmd,
		updateDagCmd,
		[]stringVar{
			stringVar{&dagName, "dagName", "", "", "DAG Name", "DAG_NAME", true},
			stringVar{&dagFileName, "dagFileName", "", "", "DAG File Name", "DAG_FILE_NAME", true},
			stringVar{&dagFilePath, "dagFilePath", "", "", "DAG File Path", "DAG_FILE_PATH", true},
			stringVar{&dagVersion, "dagVersion", "", "", "DAG Version", "DAG_VERSION", true},
			stringVar{&dagDesc, "dagDesc", "", "", "DAG Description", "DAG_DESCRIPTION", true},
			stringVar{&dagFlowType, "dagFlowType", "", "", "DAG Flow Type (SPARK_JAVA or SPARK_PYTHON)", "DAG_FLOW_TYPE", true},
			stringVar{&dagTemplate, "dagTemplate", "", "", "DAG Template", "DAG_TEMPLATE", true},
		},
		[]boolVar{},
		[]intVar{})
	// deploy
	deployDagPI = NewPI(
		dagCmd,
		deployDagCmd,
		[]stringVar{
			stringVar{&dagName, "dagName", "", "", "DAG Name", "DAG_NAME", true},
		},
		[]boolVar{},
		[]intVar{})
	// status
	dagStatusPI = NewPI(
		dagCmd,
		dagStatusCmd,
		[]stringVar{
			stringVar{&dagName, "dagName", "", "", "DAG Name", "DAG_NAME", true},
		},
		[]boolVar{},
		[]intVar{})
	// task-run-info
	getDagTaskRunPI = NewPI(
		dagCmd,
		getDagTaskRunCmd,
		[]stringVar{
			stringVar{&dagName, "dagName", "", "", "DAG Name", "DAG_NAME", true},
			stringVar{&dagRunID, "dagRunID", "", "", "DAG Run ID", "DAG_RUN_ID", true},
			stringVar{&dagRunID, "dagTaskID", "", "", "DAG Task ID", "DAG_TASK_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	// list-run
	getDagRunPI = NewPI(
		dagCmd,
		getDagRunCmd,
		[]stringVar{
			stringVar{&dagName, "dagName", "", "", "DAG Name", "DAG_NAME", true},
			stringVar{&dagRunID, "dagRunID", "", "", "DAG Run ID", "DAG_RUN_ID", false},
		},
		[]boolVar{},
		[]intVar{})
	// list-task
	getDagTaskPI = NewPI(
		dagCmd,
		getDagTaskCmd,
		[]stringVar{
			stringVar{&dagName, "dagName", "", "", "DAG Name", "DAG_NAME", true},
			stringVar{&dagRunID, "dagTaskID", "", "", "DAG Task ID", "DAG_TASK_ID", false},
		},
		[]boolVar{},
		[]intVar{})
	// DEPENDENCY Commands
	// list
	getDependencyPI = NewPI(
		dependencyCmd,
		getDependencyCmd,
		[]stringVar{
			stringVar{&dependencyID, "dependencyID", "", "", "Dependency ID", "DEPENDENCY_ID", false},
		},
		[]boolVar{},
		[]intVar{})
	// delete
	deleteDependencyPI = NewPI(
		dependencyCmd,
		deleteDependencyCmd,
		[]stringVar{
			stringVar{&dependencyID, "dependencyID", "", "", "Dependency ID", "DEPENDENCY_ID", true},
		},
		[]boolVar{
			boolVar{&force, "force", "f", false, "Permanently remove a DAG", "FORCE", false},
		},
		[]intVar{})
	// deploy
	deployDependencyPI = NewPI(
		dependencyCmd,
		deployDependencyCmd,
		[]stringVar{
			stringVar{&dependencyID, "dependencyID", "", "", "Dependency ID", "DEPENDENCY_ID", false},
		},
		[]boolVar{},
		[]intVar{})
	// undeploy
	unDeployDependencyPI = NewPI(
		dependencyCmd,
		unDeployDependencyCmd,
		[]stringVar{
			stringVar{&dependencyID, "dependencyID", "", "", "Dependency ID", "DEPENDENCY_ID", false},
		},
		[]boolVar{},
		[]intVar{})
	// create
	postDependencyPI = NewPI(
		dependencyCmd,
		postDependencyCmd,
		[]stringVar{
			stringVar{&dependencyType, "dependencyType", "", "", "Dependency Type", "DEPENDENCY_TYPE", true},
			stringVar{&dependencyFileName, "dependencyFileName", "", "", "Dependency File Name", "DEPENDENCY_FILE_NAME", true},
			stringVar{&dependencyFileLocation, "dependencyFileLocation", "", "", "Dependency File Location", "DEPENDENCY_FILE_LOCATION", true},
		},
		[]boolVar{},
		[]intVar{})

	// FLOW TEMPLATE Commands
	// create
	postFlowTemplatePI = NewPI(
		flowTemplateCmd,
		postFlowTemplateCmd,
		[]stringVar{
			stringVar{&flowTemplateName, "flowTemplateName", "", "", "Flow Template Name", "FLOW_TEMPLATE_NAME", true},
			stringVar{&templateFileName, "templateFileName", "", "", "Flow Template Analytic File Name", "TEMPLATE_FILE_NAME", true},
			stringVar{&templateFilePath, "templateFilePath", "", "", "Flow Template Analytic File Path", "TEMPLATE_FILE_PATH", true},
			stringVar{&flowTemplateVersion, "flowTemplateVersion", "", "", "Flow Template Version", "FLOW_TEMPLATE_VERSION", true},
			stringVar{&desc, "desc", "", "", "Flow Template Description", "DESC", true},
			stringVar{&flowType, "flowType", "", "", "Flow Type (SPARK_JAVA or SPARK_PYTHON)", "FLOW_TYPE", true},
		},
		[]boolVar{},
		[]intVar{})
	// update
	updateFlowTemplatePI = NewPI(
		flowTemplateCmd,
		updateFlowTemplateCmd,
		[]stringVar{
			stringVar{&flowTemplateID, "flowTemplateID", "", "", "Flow Template ID", "FLOW_TEMPLATE_ID", true},
			stringVar{&flowTemplateName, "flowTemplateName", "", "", "Flow Template Name", "FLOW_TEMPLATE_NAME", true},
			stringVar{&templateFileName, "templateFileName", "", "", "Flow Template Analytic File Name", "TEMPLATE_FILE_NAME", true},
			stringVar{&templateFilePath, "templateFilePath", "", "", "Flow Template Analytic File Path", "TEMPLATE_FILE_PATH", true},
			stringVar{&flowTemplateVersion, "flowTemplateVersion", "", "", "Flow Template Version", "FLOW_TEMPLATE_VERSION", true},
			stringVar{&desc, "desc", "", "", "Flow Template Description", "DESC", true},
			stringVar{&flowType, "flowType", "", "", "Flow Type (SPARK_JAVA or SPARK_PYTHON)", "FLOW_TYPE", true},
		},
		[]boolVar{},
		[]intVar{})
	// update-spark-args
	updateFlowTemplateChangeSparkArgumentsPI = NewPI(
		flowTemplateCmd,
		updateFlowTemplateChangeSparkArguments,
		[]stringVar{
			stringVar{&flowTemplateID, "flowTemplateID", "", "", "Flow Template ID", "FLOW_TEMPLATE_ID", true},
			stringVar{&sparkArgs, "sparkArgs", "", "", "Flow Encapsulated Spark Arguments", "SPARK_ARGS", true},
		},
		[]boolVar{},
		[]intVar{})
	// list
	getFlowTemplatePI = NewPI(
		flowTemplateCmd,
		getFlowTemplateCmd,
		[]stringVar{
			stringVar{&flowTemplateID, "flowTemplateID", "", "", "Flow Template ID", "FLOW_TEMPLATE_ID", false},
			stringVar{&flowTemplateName, "flowTemplateName", "", "", "Flow Template Name", "FLOW_TEMPLATE_NAME", false},
		},
		[]boolVar{},
		[]intVar{})
	// delete
	deleteFlowTemplatePI = NewPI(
		flowTemplateCmd,
		deleteFlowTemplateCmd,
		[]stringVar{
			stringVar{&flowTemplateID, "flowTemplateID", "", "", "Flow Template ID", "FLOW_TEMPLATE_ID", true},
		},
		[]boolVar{
			boolVar{&force, "force", "f", false, "Permanently remove a DAG", "FORCE", false},
		},
		[]intVar{})
	// list-tags
	getFlowTemplateTagsPI = NewPI(
		flowTemplateCmd,
		getFlowTemplateTagsCmd,
		[]stringVar{
			stringVar{&flowTemplateID, "flowTemplateID", "", "", "Flow Template ID", "FLOW_TEMPLATE_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	// save-tags
	saveFlowTemplateTagsPI = NewPI(
		flowTemplateCmd,
		saveFlowTemplateTagsCmd,
		[]stringVar{
			stringVar{&flowTemplateID, "flowTemplateID", "", "", "Flow Template ID", "FLOW_TEMPLATE_ID", true},
			stringVar{&tags, "tags", "", "", "Flow Template Tags", "TAGS", true},
		},
		[]boolVar{},
		[]intVar{})

	// FLOW Commands
	// list
	getFlowPI = NewPI(
		flowCmd,
		getFlowCmd,
		[]stringVar{
			stringVar{&flowName, "flowName", "", "", "Flow Name", "FLOW_NAME", false},
			stringVar{&flowID, "flowID", "", "", "Flow ID", "FLOW_ID", false},
			stringVar{&flowTemplateID, "flowTemplateID", "", "", "Flow Template ID", "FLOW_TEMPLATE_ID", false},
		},
		[]boolVar{},
		[]intVar{})
	// create
	postFlowPI = NewPI(
		flowCmd,
		postFlowCmd,
		[]stringVar{
			stringVar{&flowName, "flowName", "", "", "Flow Name", "FLOW_NAME", true},
			stringVar{&flowTemplateID, "flowTemplateID", "", "", "Flow Template ID", "FLOW_TEMPLATE_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	// create-direct
	postDirectFlowPI = NewPI(
		flowCmd,
		postDirectFlowCmd,
		[]stringVar{
			stringVar{&flowName, "flowName", "", "", "Flow Name", "FLOW_NAME", true},
			stringVar{&flowFileName, "flowFileName", "", "", "Flow File Name", "FLOW_FILE_NAME", true},
			stringVar{&flowFilePath, "flowFilePath", "", "", "Flow File Path", "FLOW_FILE_PATH", true},
			stringVar{&flowVersion, "flowVersion", "", "", "Direct Flow Version", "FLOW_VERSION", true},
			stringVar{&desc, "desc", "", "", "Flow Template Description", "DESC", true},
			stringVar{&flowType, "flowType", "", "", "Flow Type (SPARK_JAVA or SPARK_PYTHON)", "FLOW_TYPE", true},
		},
		[]boolVar{},
		[]intVar{})
	// create-flow-template
	createFlowTemplateFromFlowPI = NewPI(
		flowCmd,
		createFlowTemplateFromFlowCmd,
		[]stringVar{
			stringVar{&flowID, "flowID", "", "", "Flow ID", "FLOW_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	// add-config-file
	addFlowConfigFilesPI = NewPI(
		flowCmd,
		addFlowConfigFiles,
		[]stringVar{
			stringVar{&configFileDetails, "configFileDetails", "", "", "Flow Config File Details", "CONFIG_FILE_DETAILS", true},
			stringVar{&flowID, "flowID", "", "", "Flow ID", "FLOW_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	// list-config-file
	listConfigFilesPI = NewPI(
		flowCmd,
		listConfigFilesCmd,
		[]stringVar{
			stringVar{&flowID, "flowID", "", "", "Flow ID", "FLOW_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	// delete-config-file
	deleteFlowConfigFilePI = NewPI(
		flowCmd,
		deleteFlowConfigFile,
		[]stringVar{
			stringVar{&flowID, "flowID", "", "", "Flow ID", "FLOW_ID", true},
			stringVar{&configFileName, "configFileName", "", "", "Flow Config File Name", "CONFIG_FILE_NAME", true},
		},
		[]boolVar{},
		[]intVar{})
	// delete
	deleteFlowPI = NewPI(
		flowCmd,
		deleteFlowCmd,
		[]stringVar{
			stringVar{&flowID, "flowID", "", "", "Flow ID", "FLOW_ID", true},
		},
		[]boolVar{
			boolVar{&force, "force", "f", false, "Permanently remove a Flow", "FORCE", false},
		},
		[]intVar{})
	// update
	updateFlowChangeSparkArgumentsPI = NewPI(
		flowCmd,
		updateFlowChangeSparkArguments,
		[]stringVar{
			stringVar{&sparkArgs, "sparkArgs", "", "", "Flow Encapsulated Spark Arguments", "SPARK_ARGS", true},
			stringVar{&flowID, "flowID", "", "", "Flow ID", "FLOW_ID", true},
			stringVar{&flowTemplateID, "flowTemplateID", "", "", "Flow Template ID", "FLOW_TEMPLATE_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	// update
	updateDirectFlowPI = NewPI(
		flowCmd,
		updateDirectFlowCmd,
		[]stringVar{
			stringVar{&flowID, "flowID", "", "", "Flow ID", "FLOW_ID", true},
			stringVar{&flowFileName, "flowFileName", "", "", "Flow File Name", "FLOW_FILE_NAME", true},
			stringVar{&flowFilePath, "flowFilePath", "", "", "Flow File Path", "FLOW_FILE_PATH", true},
			stringVar{&desc, "desc", "", "", "Flow Template Description", "DESC", true},
		},
		[]boolVar{},
		[]intVar{})
	// launch
	postLaunchFlowPI = NewPI(
		flowCmd,
		postLaunchFlowCmd,
		[]stringVar{
			stringVar{&flowID, "flowID", "", "", "Flow ID", "FLOW_ID", true},
			stringVar{&flowTemplateID, "flowTemplateID", "", "", "Flow Template ID", "FLOW_TEMPLATE_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	// stop
	stopFlowPI = NewPI(
		flowCmd,
		stopFlowCmd,
		[]stringVar{
			stringVar{&flowName, "flowName", "", "", "Flow Name", "FLOW_NAME", true},
		},
		[]boolVar{},
		[]intVar{})
	// list-tags
	getFlowTagsPI = NewPI(
		flowCmd,
		getFlowTagsCmd,
		[]stringVar{
			stringVar{&flowID, "flowID", "", "", "Flow ID", "FLOW_ID", true},
			stringVar{&flowTemplateID, "flowTemplateID", "", "", "Flow Template ID", "FLOW_TEMPLATE_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	// save-tags
	saveFlowTagsPI = NewPI(
		flowCmd,
		saveFlowTagsCmd,
		[]stringVar{
			stringVar{&tags, "tags", "", "", "Flow Template Tags", "TAGS", true},
			stringVar{&flowID, "flowID", "", "", "Flow ID", "FLOW_ID", true},
			stringVar{&flowTemplateID, "flowTemplateID", "", "", "Flow Template ID", "FLOW_TEMPLATE_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	// FLOW INSTANCE Commands
	// list-instances
	getInstancePI = NewPI(
		instanceCmd,
		getInstanceCmd,
		[]stringVar{
			stringVar{&instanceID, "instanceID", "", "", "Instance ID", "INSTANCE_ID", false},
		},
		[]boolVar{},
		[]intVar{})
	// list-containers
	getAllInstanceContainersPI = NewPI(
		instanceCmd,
		getAllInstanceContainers,
		[]stringVar{
			stringVar{&instanceID, "instanceID", "", "", "Instance ID", "INSTANCE_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	// list-container-logs-response
	getContainerLogsResponsePI = NewPI(
		instanceCmd,
		getContainerLogsResponse,
		[]stringVar{
			stringVar{&instanceID, "instanceID", "", "", "Instance ID", "INSTANCE_ID", true},
			stringVar{&containerID, "containerID", "", "", "Container ID", "CONTAINER_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	// list-container-logs
	getContainerLogsPI = NewPI(
		instanceCmd,
		getContainerLogs,
		[]stringVar{
			stringVar{&instanceID, "instanceID", "", "", "Instance ID", "INSTANCE_ID", true},
			stringVar{&containerID, "containerID", "", "", "Container ID", "CONTAINER_ID", true},
		},
		[]boolVar{
			boolVar{&tail, "tail", "t", false, "Tail Container Logs", "TAIL", false},
		},
		[]intVar{
			intVar{&containerLogSink, "containerLogSink", "", 1, "Container Log Sink (stderr or stdout)", "CONTAINER_LOG_SINK", false},
		})
	// list-submit-logs
	getInstanceSubmitLogsPI = NewPI(
		instanceCmd,
		getInstanceSubmitLogsCmd,
		[]stringVar{
			stringVar{&instanceID, "instanceID", "", "", "Instance ID", "INSTANCE_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	// stop-instance
	stopInstancePI = NewPI(
		instanceCmd,
		stopInstanceCmd,
		[]stringVar{
			stringVar{&instanceID, "instanceID", "", "", "Instance ID", "INSTANCE_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	//list-spark-app-details
	getSparkAppDetailsPI = NewPI(
		instanceCmd,
		getSparkAppDetails,
		[]stringVar{
			stringVar{&instanceID, "instanceID", "", "", "Instance ID", "INSTANCE_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	//list-spark-executor-details
	getSparkExecutorDetailsPI = NewPI(
		instanceCmd,
		getSparkExecutorDetails,
		[]stringVar{
			stringVar{&instanceID, "instanceID", "", "", "Instance ID", "INSTANCE_ID", true},
			stringVar{&attemptID, "attemptID", "", "", "Attempt ID", "ATTEMPT_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	//list-app-stages
	getAllAppStagesPI = NewPI(
		instanceCmd,
		getAllAppStages,
		[]stringVar{
			stringVar{&instanceID, "instanceID", "", "", "Instance ID", "INSTANCE_ID", true},
			stringVar{&attemptID, "attemptID", "", "", "Attempt ID", "ATTEMPT_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	//list-attempts
	getAllAttemptsPI = NewPI(
		instanceCmd,
		getAllAttempts,
		[]stringVar{
			stringVar{&instanceID, "instanceID", "", "", "Instance ID", "INSTANCE_ID", true},
			stringVar{&attemptID, "attemptID", "", "", "Attempt ID", "ATTEMPT_ID", true},
			stringVar{&stageID, "stageID", "", "", "Stage ID", "STAGE_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	//list-attempt-details
	getAttemptDetailsPI = NewPI(
		instanceCmd,
		getAttemptDetails,
		[]stringVar{
			stringVar{&instanceID, "instanceID", "", "", "Instance ID", "INSTANCE_ID", true},
			stringVar{&attemptID, "attemptID", "", "", "Attempt ID", "ATTEMPT_ID", true},
			stringVar{&stageID, "stageID", "", "", "Stage ID", "STAGE_ID", true},
			stringVar{&stageAttemptID, "stageAttemptID", "", "", "Stage Attempt ID", "STAGE_ATTEMPT_ID", true},
		},
		[]boolVar{},
		[]intVar{})
	//list-tasks
	getAllTasksByStagePI = NewPI(
		instanceCmd,
		getAllTasksByStage,
		[]stringVar{
			stringVar{&instanceID, "instanceID", "", "", "Instance ID", "INSTANCE_ID", true},
			stringVar{&attemptID, "attemptID", "", "", "Attempt ID", "ATTEMPT_ID", true},
			stringVar{&stageID, "stageID", "", "", "Stage ID", "STAGE_ID", true},
			stringVar{&stageAttemptID, "stageAttemptID", "", "", "Stage Attempt ID", "STAGE_ATTEMPT_ID", true},
		},
		[]boolVar{},
		[]intVar{})

	// list all commands
	commands = []*pi{&loginPI, &healthCheckPI, &versionCheckPI, &getDagPI, &deleteDagPI, &postDagPI, &updateDagPI, &deployDagPI, &dagStatusPI, &getDagTaskRunPI, &getDagRunPI, &getDagTaskPI, &getDependencyPI, &deleteDependencyPI, &deployDependencyPI, &unDeployDependencyPI, &postDependencyPI, &postFlowTemplatePI, &updateFlowTemplatePI, &updateFlowTemplateChangeSparkArgumentsPI, &getFlowTemplatePI, &deleteFlowTemplatePI, &getFlowTemplateTagsPI, &saveFlowTemplateTagsPI, &getFlowPI, &postFlowPI, &postDirectFlowPI, &createFlowTemplateFromFlowPI, &addFlowConfigFilesPI, &listConfigFilesPI, &deleteFlowConfigFilePI, &deleteFlowPI, &updateFlowChangeSparkArgumentsPI, &updateDirectFlowPI, &postLaunchFlowPI, &stopFlowPI, &getFlowTagsPI, &saveFlowTagsPI, &getInstancePI, &getAllInstanceContainersPI, &getContainerLogsResponsePI, &getContainerLogsPI, &getInstanceSubmitLogsPI, &stopInstancePI, &getSparkAppDetailsPI, &getSparkExecutorDetailsPI, &getAllAppStagesPI, &getAllAttemptsPI, &getAttemptDetailsPI, &getAllTasksByStagePI}

	// GENERAL GLOBAL flags
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "", homeDir+"/.pi/"+file, "config file location")
	viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	viper.BindEnv("config", "CONFIG")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbosity")
	viper.BindPFlag("verbose", RootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindEnv("verbose", "VERBOSE")
	RootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "i", false, "Enable interactive mode (prompt user for input)")
	viper.BindPFlag("interactive", RootCmd.PersistentFlags().Lookup("interactive"))
	viper.BindEnv("interactive", "INTERACTIVE")

	// configure command structure
	for _, c := range commands {
		c.parent.AddCommand(c.C)
	}

	// set PI CLI version
	RootCmd.Version = Version + "\ngit commit hash " + GitHash

	// set up flags for all sub commands
	for _, c := range commands {
		for _, f := range c.strFlags {
			c.C.PersistentFlags().StringVarP(f.p, f.name, f.shorthand, f.value, f.usage)
			c.V.BindPFlag(f.name, c.C.PersistentFlags().Lookup(f.name))
			c.V.BindEnv(f.name, f.env)
		}
		for _, f := range c.boolFlags {
			c.C.PersistentFlags().BoolVarP(f.p, f.name, f.shorthand, f.value, f.usage)
			c.V.BindPFlag(f.name, c.C.PersistentFlags().Lookup(f.name))
			c.V.BindEnv(f.name, f.env)
		}
		for _, f := range c.intFlags {
			c.C.PersistentFlags().IntVarP(f.p, f.name, f.shorthand, f.value, f.usage)
			c.V.BindPFlag(f.name, c.C.PersistentFlags().Lookup(f.name))
			c.V.BindEnv(f.name, f.env)
		}
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Use config file from the flag.
	if viper.GetString("config") != "" {
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		fmt.Println(errors.New("error configuration file is set to: ''"))
		os.Exit(1)
	}

	//viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}

	// if not given via command line or env. var try to get from config file
	for _, c := range commands {
		_ = getConfigFileParams(flags, c)
		if !viper.GetBool("interactive") {
			markRequired(c)
		}
	}
}
