package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	APIHost                                  string
	TenantID                                 string
	IssuerID                                 string
	ClientID                                 string
	ClientSecret                             string
	Token                                    string
	homeDir                                  = os.Getenv("HOME")
	dir                                      = homeDir + "/.pi"
	file                                     = "config.json"
	deleteFlowTemplateID                     string
	flowTemplateID                           string
	flowTemplateName                         string
	templateFileName                         string
	templateFilePath                         string
	flowVersion                              string
	flowTemplateVersion                      string
	desc                                     string
	flowType                                 string
	flowName                                 string
	flowID                                   string
	sparkArgs                                string
	instanceID                               string
	containerID                              string
	dagName                                  string
	dagFileName                              string
	dagFilePath                              string
	dagVersion                               string
	dagDesc                                  string
	dagFlowType                              string
	dagTemplate                              string
	dagRunID                                 string
	dagTaskID                                string
	dependencyID                             string
	dependencyType                           string
	dependencyFileName                       string
	dependencyFileLocation                   string
	flowFileName                             string
	flowFilePath                             string
	configFileDetails                        string
	configFileName                           string
	tags                                     string
	attemptID                                string
	stageID                                  string
	stageAttemptID                           string
	containerLogSink                         int
	tail                                     bool
	cfgFile                                  string
	verbose                                  bool
	interactive                              bool
	force                                    bool
	Version                                  = "No Version Provided"
	GitHash                                  = "No GitHash Provided"
	GitDate                                  = "No GitDate Provided"
	loginPI                                  = pi{}
	healthCheckPI                            = pi{}
	versionCheckPI                           = pi{}
	getDagPI                                 = pi{}
	deleteDagPI                              = pi{}
	postDagPI                                = pi{}
	updateDagPI                              = pi{}
	deployDagPI                              = pi{}
	dagStatusPI                              = pi{}
	getDagTaskRunPI                          = pi{}
	getDagRunPI                              = pi{}
	getDagTaskPI                             = pi{}
	getDependencyPI                          = pi{}
	deleteDependencyPI                       = pi{}
	deployDependencyPI                       = pi{}
	unDeployDependencyPI                     = pi{}
	postDependencyPI                         = pi{}
	postFlowTemplatePI                       = pi{}
	updateFlowTemplatePI                     = pi{}
	updateFlowTemplateChangeSparkArgumentsPI = pi{}
	getFlowTemplatePI                        = pi{}
	deleteFlowTemplatePI                     = pi{}
	getFlowTemplateTagsPI                    = pi{}
	saveFlowTemplateTagsPI                   = pi{}
	getFlowPI                                = pi{}
	postFlowPI                               = pi{}
	postDirectFlowPI                         = pi{}
	createFlowTemplateFromFlowPI             = pi{}
	addFlowConfigFilesPI                     = pi{}
	listConfigFilesPI                        = pi{}
	deleteFlowConfigFilePI                   = pi{}
	deleteFlowPI                             = pi{}
	updateFlowChangeSparkArgumentsPI         = pi{}
	updateDirectFlowPI                       = pi{}
	postLaunchFlowPI                         = pi{}
	stopFlowPI                               = pi{}
	getFlowTagsPI                            = pi{}
	saveFlowTagsPI                           = pi{}
	getInstancePI                            = pi{}
	getAllInstanceContainersPI               = pi{}
	getContainerLogsResponsePI               = pi{}
	getContainerLogsPI                       = pi{}
	getInstanceSubmitLogsPI                  = pi{}
	stopInstancePI                           = pi{}
	getSparkAppDetailsPI                     = pi{}
	getSparkExecutorDetailsPI                = pi{}
	getAllAppStagesPI                        = pi{}
	getAllAttemptsPI                         = pi{}
	getAttemptDetailsPI                      = pi{}
	getAllTasksByStagePI                     = pi{}
	flags                                    = []flag{flag{"APIHost", "string"}, flag{"ClientID", "string"}, flag{"ClientSecret", "string"}, flag{"IssuerID", "string"}, flag{"TenantID", "string"}, flag{"Token", "string"}, flag{"attemptID", "string"}, flag{"configFileDetails", "string"}, flag{"configFileName", "string"}, flag{"containerID", "string"}, flag{"containerLogSink", "int"}, flag{"dagDesc", "string"}, flag{"dagFileName", "string"}, flag{"dagFilePath", "string"}, flag{"dagFlowType", "string"}, flag{"dagName", "string"}, flag{"dagRunID", "string"}, flag{"dagTaskID", "string"}, flag{"dagTemplate", "string"}, flag{"dagVersion", "string"}, flag{"dependencyFileLocation", "string"}, flag{"dependencyFileName", "string"}, flag{"dependencyID", "string"}, flag{"dependencyType", "string"}, flag{"desc", "string"}, flag{"flowFileName", "string"}, flag{"flowFilePath", "string"}, flag{"flowID", "string"}, flag{"flowName", "string"}, flag{"flowTemplateID", "string"}, flag{"flowTemplateName", "string"}, flag{"flowTemplateVersion", "string"}, flag{"flowType", "string"}, flag{"flowVersion", "string"}, flag{"force", "bool"}, flag{"instanceID", "string"}, flag{"sparkArgs", "string"}, flag{"stageAttemptID", "string"}, flag{"stageID", "string"}, flag{"tags", "string"}, flag{"tail", "bool"}, flag{"templateFileName", "string"}, flag{"templateFilePath", "string"}, flag{"verbose", "bool"}, flag{"interactive", "bool"}}
	commands                                 = []*pi{}
)

type flag struct {
	Name string
	Type string
}

type pi struct {
	parent    *cobra.Command
	C         *cobra.Command
	V         *viper.Viper
	strFlags  []stringVar
	boolFlags []boolVar
	intFlags  []intVar
}

func NewPI(parent, command *cobra.Command, strFlags []stringVar, boolFlags []boolVar, intFlags []intVar) pi {
	return pi{parent: parent, C: command, V: viper.New(), strFlags: strFlags, boolFlags: boolFlags, intFlags: intFlags}
}

func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		if err.Error() != "unexpected newline" {
			fmt.Println("Error parsing input...")
			return false
		}
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Printf("Please type yes or no and then press enter: ")
		return askForConfirmation()
	}
}

func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}

func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

func getInputString() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		return "", scanner.Err()
	}
	return scanner.Text(), nil
}
func getInputBool() (bool, error) {
	var response bool
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false, err
	}
	return response, nil
}
func getInputInt() (int, error) {
	var response int
	_, err := fmt.Scanln(&response)
	if err != nil {
		return 0, err
	}
	return response, nil
}

func getMissingRequiredParams(pi pi) error {
	if viper.GetBool("interactive") {
		for _, p := range pi.strFlags {
			if p.required {
				if pi.V.GetString(p.name) == "" {
					fmt.Printf("Enter %s: ", p.name)
				} else {
					fmt.Printf("Enter %s (%s): ", p.name, pi.V.GetString(p.name))
				}
				response, err := getInputString()
				if err != nil {
					fmt.Println("found error")
					return err
				}
				if response != "" {
					pi.V.Set(p.name, response)
				}
			}
		}
	}
	return nil
}

func getConfigFileParams(flags []flag, pi *pi) error {
	for _, f := range flags {
		switch f.Type {
		case "string":
			if pi.V.GetString(f.Name) == "" {
				pi.V.Set(f.Name, viper.GetString(f.Name))
			}
		case "bool":
			//pi.V.Set(f.Name, viper.GetBool(f.Name))
		case "int":
			//pi.V.Set(f.Name, viper.GetInt(f.Name))
		default:
			return errors.New("Flag Type not supported: " + f.Name + ", " + f.Type)

		}
	}

	return nil
}

func markRequired(pi *pi) {
	for _, f := range pi.strFlags {
		if f.required && pi.V.GetString(f.name) == "" {
			pi.C.MarkFlagRequired(f.name)
		}
	}
}

type stringVar struct {
	p         *string
	name      string
	shorthand string
	value     string
	usage     string
	env       string
	required  bool
}
type boolVar struct {
	p         *bool
	name      string
	shorthand string
	value     bool
	usage     string
	env       string
	required  bool
}
type intVar struct {
	p         *int
	name      string
	shorthand string
	value     int
	usage     string
	env       string
	required  bool
}
