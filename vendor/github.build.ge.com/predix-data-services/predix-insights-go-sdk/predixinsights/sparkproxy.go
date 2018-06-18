package predixinsights

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type ApplicationDetails struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Attempts []Attempt `json:"attempts"`
}

type Attempt struct {
	AttemptID        string `json:"attemptId"`
	StartTime        string `json:"startTime"`
	EndTime          string `json:"endTime"`
	LastUpdated      string `json:"lastUpdated"`
	Duration         int    `json:"duration"`
	SparkUser        string `json:"sparkUser"`
	Completed        bool   `json:"completed"`
	StartTimeEpoch   int64  `json:"startTimeEpoch"`
	EndTimeEpoch     int64  `json:"endTimeEpoch"`
	LastUpdatedEpoch int64  `json:"lastUpdatedEpoch"`
}

type ExecutorDetails struct {
	ID                string `json:"id"`
	HostPort          string `json:"hostPort"`
	IsActive          bool   `json:"isActive"`
	RddBlocks         int    `json:"rddBlocks"`
	MemoryUsed        int    `json:"memoryUsed"`
	DiskUsed          int    `json:"diskUsed"`
	TotalCores        int    `json:"totalCores"`
	MaxTasks          int    `json:"maxTasks"`
	ActiveTasks       int    `json:"activeTasks"`
	FailedTasks       int    `json:"failedTasks"`
	CompletedTasks    int    `json:"completedTasks"`
	TotalTasks        int    `json:"totalTasks"`
	TotalDuration     int    `json:"totalDuration"`
	TotalGCTime       int    `json:"totalGCTime"`
	TotalInputBytes   int    `json:"totalInputBytes"`
	TotalShuffleRead  int    `json:"totalShuffleRead"`
	TotalShuffleWrite int    `json:"totalShuffleWrite"`
	MaxMemory         int64  `json:"maxMemory"`
	ExecutorLogs      struct {
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
	} `json:"executorLogs"`
}
type StageInformation struct {
	Status                string              `json:"status"`
	StageID               int                 `json:"stageId"`
	AttemptID             int                 `json:"attemptId"`
	NumActiveTasks        int                 `json:"numActiveTasks"`
	NumCompleteTasks      int                 `json:"numCompleteTasks"`
	NumFailedTasks        int                 `json:"numFailedTasks"`
	ExecutorRunTime       int                 `json:"executorRunTime"`
	ExecutorCPUTime       int64               `json:"executorCpuTime"`
	SubmissionTime        string              `json:"submissionTime"`
	FirstTaskLaunchedTime string              `json:"firstTaskLaunchedTime"`
	CompletionTime        string              `json:"completionTime"`
	InputBytes            int                 `json:"inputBytes"`
	InputRecords          int                 `json:"inputRecords"`
	OutputBytes           int                 `json:"outputBytes"`
	OutputRecords         int                 `json:"outputRecords"`
	ShuffleReadBytes      int                 `json:"shuffleReadBytes"`
	ShuffleReadRecords    int                 `json:"shuffleReadRecords"`
	ShuffleWriteBytes     int                 `json:"shuffleWriteBytes"`
	ShuffleWriteRecords   int                 `json:"shuffleWriteRecords"`
	MemoryBytesSpilled    int                 `json:"memoryBytesSpilled"`
	DiskBytesSpilled      int                 `json:"diskBytesSpilled"`
	Name                  string              `json:"name"`
	Details               string              `json:"details"`
	SchedulingPool        string              `json:"schedulingPool"`
	AccumulatorUpdates    []AccumulatorUpdate `json:"accumulatorUpdates"`
}

type AccumulatorUpdate struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Update string `json:"update"`
	Value  string `json:"value"`
}

type AllAttemptsForStage struct {
	Status                string                     `json:"status"`
	StageID               int                        `json:"stageId"`
	AttemptID             int                        `json:"attemptId"`
	NumActiveTasks        int                        `json:"numActiveTasks"`
	NumCompleteTasks      int                        `json:"numCompleteTasks"`
	NumFailedTasks        int                        `json:"numFailedTasks"`
	ExecutorRunTime       int                        `json:"executorRunTime"`
	ExecutorCPUTime       int64                      `json:"executorCpuTime"`
	FirstTaskLaunchedTime string                     `json:"firstTaskLaunchedTime"`
	InputBytes            int                        `json:"inputBytes"`
	InputRecords          int                        `json:"inputRecords"`
	OutputBytes           int                        `json:"outputBytes"`
	OutputRecords         int                        `json:"outputRecords"`
	ShuffleReadBytes      int                        `json:"shuffleReadBytes"`
	ShuffleReadRecords    int                        `json:"shuffleReadRecords"`
	ShuffleWriteBytes     int                        `json:"shuffleWriteBytes"`
	ShuffleWriteRecords   int                        `json:"shuffleWriteRecords"`
	MemoryBytesSpilled    int                        `json:"memoryBytesSpilled"`
	DiskBytesSpilled      int                        `json:"diskBytesSpilled"`
	Name                  string                     `json:"name"`
	Details               string                     `json:"details"`
	SchedulingPool        string                     `json:"schedulingPool"`
	AccumulatorUpdates    []AccumulatorUpdate        `json:"accumulatorUpdates"`
	Tasks                 map[string]Task            `json:"tasks"`
	ExecutorSummaries     map[string]ExecutorSummary `json:"executorSummary"`
}

type Task struct {
	TaskID             int           `json:"taskId"`
	Index              int           `json:"index"`
	Attempt            int           `json:"attempt"`
	LaunchTime         string        `json:"launchTime"`
	ExecutorID         string        `json:"executorId"`
	Host               string        `json:"host"`
	TaskLocality       string        `json:"taskLocality"`
	Speculative        bool          `json:"speculative"`
	AccumulatorUpdates []interface{} `json:"accumulatorUpdates"`
	TaskMetrics        struct {
		ExecutorDeserializeTime    int `json:"executorDeserializeTime"`
		ExecutorDeserializeCPUTime int `json:"executorDeserializeCpuTime"`
		ExecutorRunTime            int `json:"executorRunTime"`
		ExecutorCPUTime            int `json:"executorCpuTime"`
		ResultSize                 int `json:"resultSize"`
		JvmGcTime                  int `json:"jvmGcTime"`
		ResultSerializationTime    int `json:"resultSerializationTime"`
		MemoryBytesSpilled         int `json:"memoryBytesSpilled"`
		DiskBytesSpilled           int `json:"diskBytesSpilled"`
		InputMetrics               struct {
			BytesRead   int `json:"bytesRead"`
			RecordsRead int `json:"recordsRead"`
		} `json:"inputMetrics"`
		OutputMetrics struct {
			BytesWritten   int `json:"bytesWritten"`
			RecordsWritten int `json:"recordsWritten"`
		} `json:"outputMetrics"`
		ShuffleReadMetrics struct {
			RemoteBlocksFetched int `json:"remoteBlocksFetched"`
			LocalBlocksFetched  int `json:"localBlocksFetched"`
			FetchWaitTime       int `json:"fetchWaitTime"`
			RemoteBytesRead     int `json:"remoteBytesRead"`
			LocalBytesRead      int `json:"localBytesRead"`
			RecordsRead         int `json:"recordsRead"`
		} `json:"shuffleReadMetrics"`
		ShuffleWriteMetrics struct {
			BytesWritten   int `json:"bytesWritten"`
			WriteTime      int `json:"writeTime"`
			RecordsWritten int `json:"recordsWritten"`
		} `json:"shuffleWriteMetrics"`
	} `json:"taskMetrics"`
}

type ExecutorSummary struct {
	TaskTime           int `json:"taskTime"`
	FailedTasks        int `json:"failedTasks"`
	SucceededTasks     int `json:"succeededTasks"`
	InputBytes         int `json:"inputBytes"`
	OutputBytes        int `json:"outputBytes"`
	ShuffleRead        int `json:"shuffleRead"`
	ShuffleWrite       int `json:"shuffleWrite"`
	MemoryBytesSpilled int `json:"memoryBytesSpilled"`
	DiskBytesSpilled   int `json:"diskBytesSpilled"`
}

func (ac *Client) GetSparkApplicationDetails(instanceID string) (ApplicationDetails, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/instances/%s/sparkproxy/", ac.APIHost, instanceID), nil)
	if err != nil {
		return ApplicationDetails{}, errors.Wrap(err, "[GetSparkApplicationDetails] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ApplicationDetails{}, errors.Wrap(err, "[GetSparkApplicationDetails] Failed to execute GET request")
	}

	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return ApplicationDetails{}, errors.Wrap(err, fmt.Sprintf("[GetSparkApplicationDetails] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return ApplicationDetails{}, fmt.Errorf("[GetSparkApplicationDetails] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	// Get app_id
	var applicationDetails ApplicationDetails
	err = json.NewDecoder(res.Body).Decode(&applicationDetails)
	if err != nil {
		return ApplicationDetails{}, errors.Wrap(err, fmt.Sprintf("[GetSparkApplicationDetails] Failed to decode response. Status code: %d", res.StatusCode))
	}

	return applicationDetails, nil
}

func (ac *Client) GetSparkExecutorDetails(instanceID, attemptID string) ([]ExecutorDetails, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/instances/%s/sparkproxy/%s/executors", ac.APIHost, instanceID, attemptID), nil)
	if err != nil {
		return []ExecutorDetails{}, errors.Wrap(err, "[GetSparkExecutorDetails] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []ExecutorDetails{}, errors.Wrap(err, "[GetSparkExecutorDetails] Failed to execute GET request")
	}

	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []ExecutorDetails{}, errors.Wrap(err, fmt.Sprintf("[GetSparkExecutorDetails] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return []ExecutorDetails{}, fmt.Errorf("[GetSparkExecutorDetails] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	// Get app_id
	var executorsDetails []ExecutorDetails
	err = json.NewDecoder(res.Body).Decode(&executorsDetails)
	if err != nil {
		return []ExecutorDetails{}, errors.Wrap(err, fmt.Sprintf("[GetSparkExecutorDetails] Failed to decode response. Status code: %d", res.StatusCode))
	}

	return executorsDetails, nil
}

func (ac *Client) GetAllStagesOfApplicationInstance(instanceID, attemptID string) ([]StageInformation, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/instances/%s/sparkproxy/%s/stages", ac.APIHost, instanceID, attemptID), nil)
	if err != nil {
		return []StageInformation{}, errors.Wrap(err, "[GetSparkExecutorDetails] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []StageInformation{}, errors.Wrap(err, "[GetSparkExecutorDetails] Failed to execute GET request")
	}

	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []StageInformation{}, errors.Wrap(err, fmt.Sprintf("[GetSparkExecutorDetails] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return []StageInformation{}, fmt.Errorf("[GetSparkExecutorDetails] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	// Get app_id
	var stagesInformation []StageInformation
	err = json.NewDecoder(res.Body).Decode(&stagesInformation)
	if err != nil {
		return []StageInformation{}, errors.Wrap(err, fmt.Sprintf("[GetSparkExecutorDetails] Failed to decode response. Status code: %d", res.StatusCode))
	}

	return stagesInformation, nil
}

func (ac *Client) GetAllAttemptsByStage(instanceID, attemptID, stageID string) ([]AllAttemptsForStage, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/instances/%s/sparkproxy/%s/stages/%s", ac.APIHost, instanceID, attemptID, stageID), nil)
	if err != nil {
		return []AllAttemptsForStage{}, errors.Wrap(err, "[GetAllAttemptsByStage] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []AllAttemptsForStage{}, errors.Wrap(err, "[GetAllAttemptsByStage] Failed to execute GET request")
	}

	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []AllAttemptsForStage{}, errors.Wrap(err, fmt.Sprintf("[GetAllAttemptsByStage] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return []AllAttemptsForStage{}, fmt.Errorf("[GetAllAttemptsByStage] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	// Get app_id
	var allAttemptsForStage []AllAttemptsForStage
	err = json.NewDecoder(res.Body).Decode(&allAttemptsForStage)
	if err != nil {
		return []AllAttemptsForStage{}, errors.Wrap(err, fmt.Sprintf("[GetAllAttemptsByStage] Failed to decode response. Status code: %d", res.StatusCode))
	}

	return allAttemptsForStage, nil
}

func (ac *Client) GetStageAttemptDetails(instanceID, attemptID, stageID, stageAttemptID string) (AllAttemptsForStage, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/instances/%s/sparkproxy/%s/stages/%s/%s", ac.APIHost, instanceID, attemptID, stageID, stageAttemptID), nil)
	if err != nil {
		return AllAttemptsForStage{}, errors.Wrap(err, "[GetStageAttemptDetails] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return AllAttemptsForStage{}, errors.Wrap(err, "[GetStageAttemptDetails] Failed to execute GET request")
	}

	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return AllAttemptsForStage{}, errors.Wrap(err, fmt.Sprintf("[GetStageAttemptDetails] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return AllAttemptsForStage{}, fmt.Errorf("[GetStageAttemptDetails] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	// Get app_id
	var allAttemptsForStage AllAttemptsForStage
	err = json.NewDecoder(res.Body).Decode(&allAttemptsForStage)
	if err != nil {
		return AllAttemptsForStage{}, errors.Wrap(err, fmt.Sprintf("[GetStageAttemptDetails] Failed to decode response. Status code: %d", res.StatusCode))
	}

	return allAttemptsForStage, nil
}

func (ac *Client) GetAllTasksByStage(instanceID, attemptID, stageID, stageAttemptID string) ([]Task, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/instances/%s/sparkproxy/%s/stages/%s/%s/taskList", ac.APIHost, instanceID, attemptID, stageID, stageAttemptID), nil)
	if err != nil {
		return []Task{}, errors.Wrap(err, "[GetAllTasksByStage] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []Task{}, errors.Wrap(err, "[GetAllTasksByStage] Failed to execute GET request")
	}

	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []Task{}, errors.Wrap(err, fmt.Sprintf("[GetAllTasksByStage] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return []Task{}, fmt.Errorf("[GetAllTasksByStage] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	// Get app_id
	var tasks []Task
	err = json.NewDecoder(res.Body).Decode(&tasks)
	if err != nil {
		return []Task{}, errors.Wrap(err, fmt.Sprintf("[GetAllTasksByStage] Failed to decode response. Status code: %d", res.StatusCode))
	}

	return tasks, nil
}
