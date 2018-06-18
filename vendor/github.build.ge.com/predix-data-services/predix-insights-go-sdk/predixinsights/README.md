# predixinsights
--
    import "predix-insights-go-sdk/predixinsights"


## Usage

```go
var ErrResourceAlreadyExists = errors.New("Resource already exists")
```
ErrResourceAlreadyExists error to be thrown if resource already exists

#### type AllFlowTemplatesResponse

```go
type AllFlowTemplatesResponse struct {
	Content []struct {
		ID          string        `json:"id"`
		Created     int64         `json:"created"`
		Updated     int64         `json:"updated"`
		Version     string        `json:"version"`
		User        string        `json:"user"`
		Name        string        `json:"name"`
		Description string        `json:"description"`
		Type        string        `json:"type"`
		Tags        []interface{} `json:"tags"`
		BlobPath    string        `json:"blobPath"`
		Flows       []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"flows"`
		SparkArguments struct {
			ApplicationArgs []string `json:"applicationArgs"`
			ClassName       string   `json:"className"`
		} `json:"sparkArguments,omitempty"`
	} `json:"content"`
	Last             bool        `json:"last"`
	TotalElements    int         `json:"totalElements"`
	TotalPages       int         `json:"totalPages"`
	First            bool        `json:"first"`
	Sort             interface{} `json:"sort"`
	NumberOfElements int         `json:"numberOfElements"`
	Size             int         `json:"size"`
	Number           int         `json:"number"`
}
```


#### type ArgsRequest

```go
type ArgsRequest struct {
	Name           string                 `json:"name"`
	SparkArguments map[string]interface{} `json:"sparkArguments"`
}
```

ArgsRequest struct represents arguments for spark job

#### type Client

```go
type Client struct {
	APIHost      string
	TenantID     string
	IssuerID     string
	ClientID     string
	ClientSecret string
}
```

Client struct contains Client information

#### func  NewClient

```go
func NewClient(APIHost, TenantID, IssuerID, ClientID, ClientSecret string) *Client
```
NewClient Method to retrieve new client object

#### func (*Client) CheckStatus

```go
func (ac *Client) CheckStatus() error
```
CheckStatus Method to check status of predix insights

#### func (*Client) DeleteDAG

```go
func (ac *Client) DeleteDAG(name string) error
```
DeleteDAG Method to delete DAG

#### func (*Client) DeleteDependencyByID

```go
func (ac *Client) DeleteDependencyByID(dependencyID string) error
```
DeleteDependencyByID Method to delete dependency by ID

#### func (*Client) DeleteFlow

```go
func (ac *Client) DeleteFlow(flowTemplateID, flowID string) error
```
DeleteFlow Method to delete flow by templateID and flowID

#### func (*Client) DeleteFlowTemplate

```go
func (ac *Client) DeleteFlowTemplate(flowTemplateID string) error
```
DeleteFlowTemplate Method to delete flowTemplate by ID

#### func (*Client) DeployAllDependencies

```go
func (ac *Client) DeployAllDependencies() error
```
DeployAllDependencies Method to deploy all dependencies

#### func (*Client) DeployDAG

```go
func (ac *Client) DeployDAG(name string) error
```
DeployDAG Method to deploy DAG by name

#### func (*Client) DeployDependencyByDependencyID

```go
func (ac *Client) DeployDependencyByDependencyID(dependencyID string) error
```
DeployDependencyByDependencyID Method to deploy dependency by ID

#### func (*Client) GetAllDAGs

```go
func (ac *Client) GetAllDAGs() (GetAllDAGsResponse, error)
```
GetAllDAGs Method to retrieve all DAGs

#### func (*Client) GetAllDependencies

```go
func (ac *Client) GetAllDependencies() (DependenciesResponse, error)
```
GetAllDependencies Method to retrieve all dependencies

#### func (*Client) GetAllFlowTemplates

```go
func (ac *Client) GetAllFlowTemplates() (AllFlowTemplatesResponse, error)
```

#### func (*Client) GetAllFlowTemplatesByPage

```go
func (ac *Client) GetAllFlowTemplatesByPage(maxPages int) ([]FlowTemplate, error)
```
GetAllFlowTemplates Method to retrieve all flowTemplates

#### func (*Client) GetAllFlows

```go
func (ac *Client) GetAllFlows(maxPages int) ([]Flow, error)
```
GetAllFlows Method to retrieve flows by number of pages

#### func (*Client) GetAllFlowsByTemplateID

```go
func (ac *Client) GetAllFlowsByTemplateID(templateID string) (GetAllFlowsByTemplateIDResponse, error)
```
GetAllFlowsByTemplateID Method to get all flows for a particular template by
templateID

#### func (*Client) GetAllInstanceContainers

```go
func (ac *Client) GetAllInstanceContainers(instanceID string) ([]ContainerResponse, error)
```
GetAllInstanceContainers Method to retrieve all containers for a particular
instance by instanceID

#### func (*Client) GetAllInstances

```go
func (ac *Client) GetAllInstances() (GetAllInstancesResponse, error)
```
GetAllInstances Method to retrieve all instances

#### func (*Client) GetContainerLogsByInstanceIDAndContainerID

```go
func (ac *Client) GetContainerLogsByInstanceIDAndContainerID(instanceID, containerID string) (GetContainerLogsResponse, error)
```
GetContainerLogsByInstanceIDAndContainerID Method to retrieve container logs by
instanceID and containerID

#### func (*Client) GetDAG

```go
func (ac *Client) GetDAG(name string) (DAGResponse, error)
```
GetDAG Method to retrieve DAG by name

#### func (*Client) GetDependencyByID

```go
func (ac *Client) GetDependencyByID(dependencyID string) (DependencyResponse, error)
```
GetDependencyByID Method to retrieve dependency by ID

#### func (*Client) GetFlow

```go
func (ac *Client) GetFlow(flowName string) (Flow, error)
```
GetFlow Method to get flow by flowName

#### func (*Client) GetFlowByTemplateIDAndFlowID

```go
func (ac *Client) GetFlowByTemplateIDAndFlowID(templateID string, flowID string) (FlowResponse, error)
```
GetFlowByTemplateIDAndFlowID Method to retrieve flow by templateid and flowid

#### func (*Client) GetFlowTemplate

```go
func (ac *Client) GetFlowTemplate(flowTemplateID string) (FlowTemplate, error)
```
GetFlowTemplate Method to get flowTemplate by flowTemplateID

#### func (*Client) GetInstance

```go
func (ac *Client) GetInstance(instanceID string) (InstanceResponse, error)
```
GetInstance Method to retrieve instance by instanceID

#### func (*Client) GetInstanceContainerLogs

```go
func (ac *Client) GetInstanceContainerLogs(instanceID, containerID string, containerLogSink ContainerLogSink) (string, error)
```
GetInstanceContainerLogs Method to retrieve container error logs and stdout logs
by instanceID, containerID and containerSink(stderr,stdout)

#### func (*Client) GetInstanceSubmitLogsByInstanceID

```go
func (ac *Client) GetInstanceSubmitLogsByInstanceID(instanceID string) (string, error)
```
GetInstanceSubmitLogsByInstanceID Method to retrieve instance submit logs

#### func (*Client) GetTagsByFlowTemplateID

```go
func (ac *Client) GetTagsByFlowTemplateID(flowTemplateID string) (TagsArray, error)
```
GetTagsByFlowTemplateID Method to get all tags for flowTemplate by
flowTemplateID

#### func (*Client) GetTagsForFlowByFlowTemplateIDAndFlowID

```go
func (ac *Client) GetTagsForFlowByFlowTemplateIDAndFlowID(flowTemplateID string, flowID string) (TagsArray, error)
```
GetTagsForFlowByFlowTemplateIDAndFlowID Method to get tags for a particular flow

#### func (*Client) LaunchFlow

```go
func (ac *Client) LaunchFlow(flowTemplateID, flowID string) (LaunchResponse, error)
```
LaunchFlow Method to launch flow by templateID and flowID

#### func (*Client) PostArguments

```go
func (ac *Client) PostArguments(flowName string, flowTemplateID string, flowID string, sparkArgs map[string]interface{}) error
```
PostArguments Method to post spark arguments

#### func (*Client) PostDAG

```go
func (ac *Client) PostDAG(dagName, dagFileName, dagFilePath, version, desc, flowType string, dt DAGTemplate) (DAGResponse, error)
```
PostDAG Method to post a new DAG

#### func (*Client) PostDependency

```go
func (ac *Client) PostDependency(dependencyType, dependencyFileName, dependencyFileLocation string) ([]DependencyResponse, error)
```
PostDependency Method to post new dependency

#### func (*Client) PostFlow

```go
func (ac *Client) PostFlow(flowName, flowTemplateID string) (Flow, error)
```
PostFlow Method to post flow

#### func (*Client) PostFlowTemplate

```go
func (ac *Client) PostFlowTemplate(flowTemplateName, templateFileName, templateFilePath, version, desc, flowType string) (FlowTemplate, error)
```
PostFlowTemplate Method to post new Template

#### func (*Client) PostMultipleDependencies

```go
func (ac *Client) PostMultipleDependencies(dependencies []DependencyDetails) ([]DependencyResponse, error)
```
PostMultipleDependencies Method to post multiple dependencies

#### func (*Client) RefreshAuthToken

```go
func (ac *Client) RefreshAuthToken() error
```
RefreshAuthToken Method to refresh UAA token

#### func (*Client) RefreshCookie

```go
func (ac *Client) RefreshCookie() error
```
RefreshCookie Method to refresh cookie

#### func (*Client) SaveTagsForFlow

```go
func (ac *Client) SaveTagsForFlow(flowTemplateID string, flowID string, tagsarray TagsArray) (FlowResponse, error)
```
SaveTagsForFlow Method to save tags for a particular flow

#### func (*Client) SaveTagsForFlowTemplate

```go
func (ac *Client) SaveTagsForFlowTemplate(flowTemplateID string, tagsarray TagsArray) (SaveTagsForFlowTemplateResponse, error)
```
SaveTagsForFlowTemplate Method to save tags for flowTemplate

#### func (*Client) StopFlow

```go
func (ac *Client) StopFlow(flowName string) error
```
StopFlow Method to stop flow by flowName

#### func (*Client) StopInstance

```go
func (ac *Client) StopInstance(instanceID string) error
```
StopInstance Method to stop a particular instance by instanceID

#### func (*Client) UnDeployAllDependencies

```go
func (ac *Client) UnDeployAllDependencies() error
```
UnDeployAllDependencies Method to undeploy all dependencies

#### func (*Client) UnDeployDependencyByDependencyID

```go
func (ac *Client) UnDeployDependencyByDependencyID(dependencyID string) error
```
UnDeployDependencyByDependencyID Method to undeploy dependency by ID

#### func (*Client) UpdateDAG

```go
func (ac *Client) UpdateDAG(dagName, dagFileName, dagFilePath, version, desc, flowType string, dt DAGTemplate) error
```
UpdateDAG Method to update existing DAG

#### func (*Client) UpdateFlowByFlowTemplateIDAndFlowIDAddConfigFile

```go
func (ac *Client) UpdateFlowByFlowTemplateIDAndFlowIDAddConfigFile(flowTemplateID, flowID string, fileDetails []FileDetails) error
```
UpdateFlowByFlowTemplateIDAndFlowIDAddConfigFile Method to update flow by adding
config files

#### func (*Client) UpdateFlowChangeSparkArguments

```go
func (ac *Client) UpdateFlowChangeSparkArguments(flowTemplateID string, flowID string, encapsulatedsparkargs EncapsulatedSparkArgs) error
```
UpdateFlowChangeSparkArguments Method to update flow by changing spark arguments

#### func (*Client) UpdateFlowTemplateByFlowTemplateIDChangeSparkArguments

```go
func (ac *Client) UpdateFlowTemplateByFlowTemplateIDChangeSparkArguments(flowTemplateID string, encapsulatedsparkargs EncapsulatedSparkArgs) error
```
UpdateFlowTemplateByFlowTemplateIDChangeSparkArguments Method to update existing
flowTemplate by changing spark arguments

#### func (*Client) UpdateFlowTemplateByFlowTemplateIDUsingNewZip

```go
func (ac *Client) UpdateFlowTemplateByFlowTemplateIDUsingNewZip(flowTemplateID, flowTemplateName, templateFileName, templateFilePath, version, desc, flowType string) error
```
UpdateFlowTemplateByFlowTemplateIDUsingNewZip Method to update existing
flowTemplate by uploading new zip

#### type ContainerLogSink

```go
type ContainerLogSink int
```

ContainerLogSink int represents logSink stderr or stdout

#### type ContainerResponse

```go
type ContainerResponse struct {
	ContainerID string      `json:"containerId"`
	StartTime   int         `json:"startTime"`
	FinishTime  int         `json:"finishTime"`
	Node        interface{} `json:"node"`
	MemoryMB    int         `json:"memoryMB"`
	Vcores      int         `json:"vcores"`
}
```

ContainerResponse struct represents container information

#### type Content

```go
type Content struct {
	ID              string        `json:"id"`
	ApplicationType string        `json:"applicationType"`
	StartTime       int64         `json:"startTime"`
	FinishTime      int64         `json:"finishTime"`
	Status          string        `json:"status"`
	Tags            []interface{} `json:"tags"`
	Flow            struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"flow"`
	SubmitDetails SubmitDetail `json:"submitDetails"`
	User          string       `json:"user"`
	Name          string       `json:"name"`
}
```

Content struct represents Instance information

#### type DAG

```go
type DAG struct {
	ID       string `json:"id"`
	Created  int64  `json:"created"`
	Updated  int64  `json:"updated"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Deployed bool   `json:"deployed"`
}
```

DAG Struct represents an airflow DAG

#### type DAGResponse

```go
type DAGResponse struct {
	ID          string        `json:"id"`
	Created     int64         `json:"created"`
	Updated     int64         `json:"updated"`
	Name        string        `json:"name"`
	Description interface{}   `json:"description"`
	Type        string        `json:"type"`
	Tags        []interface{} `json:"tags"`
	BlobPath    string        `json:"blobPath"`
	Deployed    bool          `json:"deployed"`
}
```

DAGResponse struct represents a response struct containing Airflow DAG
information

#### type DAGTemplate

```go
type DAGTemplate struct {
	Owner    string
	FlowName string
	Interval string
}
```

DAGTemplate struct

#### type DependenciesResponse

```go
type DependenciesResponse struct {
	Content          []DependencyResponse `json:"content"`
	Last             bool                 `json:"last"`
	TotalPages       int                  `json:"totalPages"`
	TotalElements    int                  `json:"totalElements"`
	First            bool                 `json:"first"`
	Sort             interface{}          `json:"sort"`
	NumberOfElements int                  `json:"numberOfElements"`
	Size             int                  `json:"size"`
	Number           int                  `json:"number"`
}
```

DependenciesResponse struct representing list of Dependency and some other meta
information

#### type DependencyDetails

```go
type DependencyDetails struct {
	Type         string
	FileName     string
	FileLocation string
}
```

DependencyDetails struct representing Dependency related information

#### type DependencyResponse

```go
type DependencyResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Tenant   string `json:"tenant"`
	Type     string `json:"type"`
	Deployed bool   `json:"deployed"`
}
```

DependencyResponse struct representing Dependency information

#### type EncapsulatedSparkArgs

```go
type EncapsulatedSparkArgs struct {
	SparkArgs SparkArguments `json:"sparkArguments"`
}
```

EncapsulatedSparkArgs struct encapsulates spark arguments

#### type Environment

```go
type Environment struct {
	HADOOPCONFDIR string `json:"HADOOP_CONF_DIR"`
	SPARKHOME     string `json:"SPARK_HOME"`
	SPARKCONFDIR  string `json:"SPARK_CONF_DIR"`
}
```

Environment struct represents environment details

#### type FileDetails

```go
type FileDetails struct {
	FileName     string
	FileLocation string
	Fields       []string
	Values       []string
}
```

FileDetails struct representing File related information

#### type Flow

```go
type Flow struct {
	ID                    string         `json:"id"`
	Created               int64          `json:"created"`
	Updated               int64          `json:"updated"`
	Version               string         `json:"version"`
	Name                  string         `json:"name"`
	Description           interface{}    `json:"description"`
	Type                  string         `json:"type"`
	Tags                  []interface{}  `json:"tags"`
	SparkArgs             SparkArguments `json:"sparkArguments"`
	LatestInstanceDetails struct {
		Summary          Summary     `json:"summary"`
		Details          interface{} `json:"details"`
		FrameworkDetails interface{} `json:"frameworkDetails"`
	} `json:"latestInstanceDetails"`
	FlowTemplate struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"flowTemplate"`
}
```

Flow struct represents Flow information

#### type FlowRequest

```go
type FlowRequest struct {
	Name string `json:"name"`
}
```

FlowRequest struct contains name of flow used for making particular flow request

#### type FlowResponse

```go
type FlowResponse struct {
	ID          string         `json:"id"`
	Created     int64          `json:"created"`
	Updated     int64          `json:"updated"`
	Version     string         `json:"version"`
	Name        string         `json:"name"`
	Description interface{}    `json:"description"`
	Type        string         `json:"type"`
	Tags        []interface{}  `json:"tags"`
	SparkArgs   SparkArguments `json:"sparkArguments"`
	FlowTemp    FlowTemplate   `json:"flowTemplate"`
}
```

FlowResponse struct represents flow information

#### type FlowTemplate

```go
type FlowTemplate struct {
	ID          string         `json:"id"`
	Created     int64          `json:"created"`
	Updated     int64          `json:"updated"`
	Version     string         `json:"version"`
	User        string         `json:"user"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        string         `json:"type"`
	Tags        []string       `json:"tags"`
	SparkArgs   SparkArguments `json:"sparkArguments"`
	Flows       []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"flows"`
}
```

FlowTemplate struct represents flow template information

#### type FlowTemplatesResponse

```go
type FlowTemplatesResponse struct {
	Content []FlowTemplate `json:"content"`
}
```

FlowTemplatesResponse struct represents list of FlowTemplates

#### type FlowsResponse

```go
type FlowsResponse struct {
	Content []Flow `json:"content"`
}
```

FlowsResponse struct represents list of flow

#### type GetAllDAGsResponse

```go
type GetAllDAGsResponse struct {
	Content          []DAG `json:"content"`
	Last             bool  `json:"last"`
	Totalpages       int   `json:"totalPages"`
	Totalelements    int   `json:"totalElements"`
	First            bool  `json:"first"`
	Numberofelements int   `json:"numberOfElements"`
	Size             int   `json:"size"`
	Number           int   `json:"number"`
}
```

GetAllDAGsResponse Struct represents list of DAGs

#### type GetAllFlowsByTemplateIDResponse

```go
type GetAllFlowsByTemplateIDResponse struct {
	Content          []FlowResponse `json:"content"`
	Last             bool           `json:"last"`
	TotalElements    int            `json:"totalElements"`
	TotalPages       int            `json:"totalPages"`
	Sort             interface{}    `json:"sort"`
	First            bool           `json:"first"`
	NumberOfElements int            `json:"numberOfElements"`
	Size             int            `json:"size"`
	Number           int            `json:"number"`
}
```

GetAllFlowsByTemplateIDResponse struct represents list of flows along iwth some
meta information

#### type GetAllInstancesResponse

```go
type GetAllInstancesResponse struct {
	Contents         []Content   `json:"content"`
	Last             bool        `json:"last"`
	TotalElements    int         `json:"totalElements"`
	TotalPages       int         `json:"totalPages"`
	First            bool        `json:"first"`
	Sort             interface{} `json:"sort"`
	NumberOfElements int         `json:"numberOfElements"`
	Size             int         `json:"size"`
	Number           int         `json:"number"`
}
```

GetAllInstancesResponse struct represents list of instances with their
information along with other meta information

#### type GetContainerLogsResponse

```go
type GetContainerLogsResponse struct {
	Stdout int `json:"stdout"`
	Stderr int `json:"stderr"`
}
```

GetContainerLogsResponse struct represents logs of a container

#### type InstanceResponse

```go
type InstanceResponse struct {
	Summary struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"summary"`
}
```

InstanceResponse struct represents summary of instance

#### type LatestInstanceDetails

```go
type LatestInstanceDetails struct {
	Summary struct {
		Status    string `json:"status"`
		StartTime int64  `json:"startTime"`
	} `json:"summary"`
}
```

LatestInstanceDetails struct contains summary of latest instance of flow

#### type LaunchResponse

```go
type LaunchResponse struct {
	ID              string        `json:"id"`
	ApplicationType string        `json:"applicationType"`
	StartTime       int64         `json:"startTime"`
	FinishTime      int           `json:"finishTime"`
	Status          string        `json:"status"`
	Tags            []interface{} `json:"tags"`
	Flow            struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"flow"`
	SubmitDetails struct {
		Command     string `json:"command"`
		Environment struct {
			HADOOPCONFDIR string `json:"HADOOP_CONF_DIR"`
			SPARKHOME     string `json:"SPARK_HOME"`
			SPARKCONFDIR  string `json:"SPARK_CONF_DIR"`
		} `json:"environment"`
	} `json:"submitDetails"`
	User string `json:"user"`
	Name string `json:"name"`
}
```

LaunchResponse struct represents response of flow launch

#### type SaveTagsForFlowTemplateResponse

```go
type SaveTagsForFlowTemplateResponse struct {
	ID          string   `json:"id"`
	Created     int64    `json:"created"`
	Updated     int64    `json:"updated"`
	Version     string   `json:"version"`
	User        string   `json:"user"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Tags        []string `json:"tags"`
	Flows       []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"flows"`
}
```

SaveTagsForFlowTemplateResponse struct represents response obtained when tags
are updated

#### type SparkArguments

```go
type SparkArguments struct {
	ApplicationArgs     []string          `json:"applicationArgs"`
	ClassName           string            `json:"className"`
	DriverCores         int               `json:"driverCores"`
	DriverMemory        string            `json:"driverMemory"`
	ExecutorMemory      string            `json:"executorMemory"`
	NumExecutors        int               `json:"numExecutors"`
	SparkListeners      []string          `json:"sparkListeners"`
	LogConf             string            `json:"logConf"`
	ExecutorEnv         map[string]string `json:"executorEnv"`
	DriverJavaOptions   string            `json:"driverJavaOptions"`
	ExecutorJavaOptions string            `json:"executorJavaOptions"`
	Confs               map[string]string `json:"confs"`
	SystemProps         map[string]string `json:"systemProps"`
	FileName            string            `json:"fileName"`
}
```

SparkArguments struct represents spark arguments

#### type SubmitDetail

```go
type SubmitDetail struct {
	Command      string      `json:"command"`
	Environments Environment `json:"environment"`
}
```

SubmitDetail struct represents command information and list of environment

#### type SubmitDetails

```go
type SubmitDetails struct {
	Command     string `json:"command"`
	Environment struct {
		HADOOPCONFDIR string `json:"HADOOP_CONF_DIR"`
		SPARKHOME     string `json:"SPARK_HOME"`
		SPARKCONFDIR  string `json:"SPARK_CONF_DIR"`
	} `json:"environment"`
}
```

SubmitDetails struct represents submit command used and the hadoop environment
information for a instance

#### type Summary

```go
type Summary struct {
	ID              string        `json:"id"`
	ApplicationType string        `json:"applicationType"`
	StartTime       int64         `json:"startTime"`
	FinishTime      int64         `json:"finishTime"`
	Status          string        `json:"status"`
	Tags            []interface{} `json:"tags"`
	Flow            struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"flow"`
	SubmitDetails SubmitDetails `json:"submitDetails"`
	User          string        `json:"user"`
	Name          string        `json:"name"`
}
```

Summary struct represents flow instance summary information

#### type TagsArray

```go
type TagsArray []string
```

TagsArray []string represents list of tags

#### type UAAResponse

```go
type UAAResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	Jti         string `json:"jti"`
}
```

UAAResponse struct contains UAA token information
