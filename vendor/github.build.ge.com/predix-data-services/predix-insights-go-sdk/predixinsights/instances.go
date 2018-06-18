package predixinsights

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// InstanceResponse struct represents summary of instance
type InstanceResponse struct {
	Summary struct {
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
	} `json:"summary"`
	Details struct {
		User                           string `json:"user"`
		Progress                       float64    `json:"progress"`
		Queue                          string `json:"queue"`
		StartTime                      int64  `json:"startTime"`
		ApplicationType                string `json:"applicationType"`
		YarnApplicationState           string `json:"yarnApplicationState"`
		FinishTime                     int64  `json:"finishTime"`
		CurrentApplicationAttemptID    string `json:"currentApplicationAttemptId"`
		FinalApplicationStatus         string `json:"finalApplicationStatus"`
		RPCPort                        int    `json:"rpcPort"`
		Diagnostics                    string `json:"diagnostics"`
		OriginalTrackingURL            string `json:"originalTrackingUrl"`
		ApplicationResourceUsageReport struct {
			NumUsedContainers     int `json:"numUsedContainers"`
			NumReservedContainers int `json:"numReservedContainers"`
			UsedResources         struct {
				Memory       int `json:"memory"`
				VirtualCores int `json:"virtualCores"`
			} `json:"usedResources"`
			ReservedResources struct {
				Memory       int `json:"memory"`
				VirtualCores int `json:"virtualCores"`
			} `json:"reservedResources"`
			NeededResources struct {
				Memory       int `json:"memory"`
				VirtualCores int `json:"virtualCores"`
			} `json:"neededResources"`
			MemorySeconds int `json:"memorySeconds"`
			VcoreSeconds  int `json:"vcoreSeconds"`
		} `json:"applicationResourceUsageReport"`
		TrackingURL string `json:"trackingUrl"`
		Name        string `json:"name"`
		ID          string `json:"id"`
		Host        string `json:"host"`
	} `json:"details"`
	FrameworkDetails interface{} `json:"frameworkDetails"`
}

// ContainerResponse struct represents container information
type ContainerResponse struct {
	ContainerID string      `json:"containerId"`
	StartTime   int         `json:"startTime"`
	FinishTime  int         `json:"finishTime"`
	Node        interface{} `json:"node"`
	MemoryMB    int         `json:"memoryMB"`
	Vcores      int         `json:"vcores"`
}

// Environment struct represents environment details
type Environment struct {
	HADOOPCONFDIR string `json:"HADOOP_CONF_DIR"`
	SPARKHOME     string `json:"SPARK_HOME"`
	SPARKCONFDIR  string `json:"SPARK_CONF_DIR"`
}

// SubmitDetail struct represents command information and list of environment
type SubmitDetail struct {
	Command      string      `json:"command"`
	Environments Environment `json:"environment"`
}

// Content struct represents Instance information
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

// GetAllInstancesResponse struct represents list of instances with their information along with other meta information
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

// GetContainerLogsResponse struct represents logs of a container
type GetContainerLogsResponse struct {
	Stdout int `json:"stdout"`
	Stderr int `json:"stderr"`
}

// ContainerLogSink int represents logSink stderr or stdout
type ContainerLogSink int

const (
	stderrSink ContainerLogSink = iota
	stdoutSink
)

// GetInstance Method to retrieve instance by instanceID
func (ac *Client) GetInstance(instanceID string) (InstanceResponse, error) {

	// Create new LaunchFlow request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/instances/%s", ac.APIHost, instanceID), nil)
	if err != nil {
		return InstanceResponse{}, errors.Wrap(err, "[GetInstance] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return InstanceResponse{}, errors.Wrap(err, "[GetInstance] Failed to execute GET request")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return InstanceResponse{}, errors.Wrap(err, fmt.Sprintf("[GetInstance] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return InstanceResponse{}, fmt.Errorf("[GetInstance] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	// Get app_id
	var instanceResponse InstanceResponse
	err = json.NewDecoder(res.Body).Decode(&instanceResponse)
	if err != nil {
		return InstanceResponse{}, errors.Wrap(err, fmt.Sprintf("[GetInstance] Failed to decode app response. Status code: %d", res.StatusCode))
	}

	return instanceResponse, nil
}

// GetAllInstances Method to retrieve all instances
func (ac *Client) GetAllInstances() (GetAllInstancesResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/instances", ac.APIHost), nil)
	if err != nil {
		return GetAllInstancesResponse{}, errors.Wrap(err, "[GetAllInstances] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return GetAllInstancesResponse{}, errors.Wrap(err, "[GetAllInstances] Failed to execute GET request")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return GetAllInstancesResponse{}, errors.Wrap(err, fmt.Sprintf("[GetAllInstances] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return GetAllInstancesResponse{}, fmt.Errorf("[GetAllInstances] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	// Get app_id
	var getAllInstancesResponse GetAllInstancesResponse
	err = json.NewDecoder(res.Body).Decode(&getAllInstancesResponse)
	if err != nil {
		return GetAllInstancesResponse{}, errors.Wrap(err, fmt.Sprintf("[GetAllInstances] Failed to decode response. Status code: %d", res.StatusCode))
	}

	return getAllInstancesResponse, nil

}

// GetAllInstanceContainers Method to retrieve all containers for a particular instance by instanceID
func (ac *Client) GetAllInstanceContainers(instanceID string) ([]ContainerResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/instances/%s/containers/", ac.APIHost, instanceID), nil)
	if err != nil {
		return []ContainerResponse{}, errors.Wrap(err, "[GetAllInstanceContainers] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []ContainerResponse{}, errors.Wrap(err, "[GetAllInstanceContainers] Failed to execute GET request")
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []ContainerResponse{}, errors.Wrap(err, fmt.Sprintf("[GetAllInstanceContainersByInstanceId] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return []ContainerResponse{}, fmt.Errorf("[GetAllInstanceContainersByInstanceId] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	// Get app_id
	var containersResponse []ContainerResponse
	err = json.NewDecoder(res.Body).Decode(&containersResponse)
	if err != nil {
		return []ContainerResponse{}, errors.Wrap(err, fmt.Sprintf("[GetAllInstanceContainersByInstanceId] Failed to decode response. Status code: %d", res.StatusCode))
	}

	return containersResponse, nil
}

// StopInstance Method to stop a particular instance by instanceID
func (ac *Client) StopInstance(instanceID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/instances/%s", ac.APIHost, instanceID), nil)
	if err != nil {
		return errors.Wrap(err, "[StopInstance] Failed to create DELETE request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[StopInstance] Failed to execute DELETE request")
	}

	defer res.Body.Close()

	// Status code is not in range of 20X
	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[StopRunningInstance] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[StopRunningInstance] Request returned %d. Body: %s", res.StatusCode, string(body))
	}
	return nil
}

// GetContainerLogsByInstanceIDAndContainerID Method to retrieve container logs by instanceID and containerID
func (ac *Client) GetContainerLogsByInstanceIDAndContainerID(instanceID, containerID string) (GetContainerLogsResponse, error) {
	// Create new LaunchFlow request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/instances/%s/containers/%s/logs", ac.APIHost, instanceID, containerID), nil)
	if err != nil {
		return GetContainerLogsResponse{}, errors.Wrap(err, "[GetContainerLogsByInstanceIDAndContainerID] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return GetContainerLogsResponse{}, errors.Wrap(err, "[GetContainerLogsByInstanceIDAndContainerID] Failed to execute GET request")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return GetContainerLogsResponse{}, errors.Wrap(err, fmt.Sprintf("[GetContainerLogsByInstanceIDAndContainerID] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return GetContainerLogsResponse{}, fmt.Errorf("[GetContainerLogsByInstanceIDAndContainerID] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	var getContainerLogsResponse GetContainerLogsResponse
	err = json.NewDecoder(res.Body).Decode(&getContainerLogsResponse)
	if err != nil {
		return GetContainerLogsResponse{}, errors.Wrap(err, fmt.Sprintf("[GetContainerLogsByInstanceIDAndContainerID] Failed to decode response. Status code: %d", res.StatusCode))
	}

	return getContainerLogsResponse, nil
}

// GetInstanceContainerLogs Method to retrieve container error logs and stdout logs by instanceID, containerID and containerSink(stderr,stdout)
func (ac *Client) GetInstanceContainerLogs(instanceID, containerID string, containerLogSink ContainerLogSink) (string, error) {
	var req *http.Request
	_ = req
	var err error
	_ = err

	switch containerLogSink {
	case stderrSink:
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/api/v1/instances/%s/containers/%s/logs/stderr", ac.APIHost, instanceID, containerID), nil)
	case stdoutSink:
		req, err = http.NewRequest("GET", fmt.Sprintf("%s/api/v1/instances/%s/containers/%s/logs/stdout", ac.APIHost, instanceID, containerID), nil)

	default:
		return "", fmt.Errorf("[GetInstanceContainerLogs] Request failed. Invalid ContainerLogSink provided")
	}

	if err != nil {
		return "", errors.Wrap(err, "[GetInstanceContainerLogs] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "[GetInstanceContainerLogs] Failed to execute GET request")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("[GetInstanceContainerLogs] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return "", fmt.Errorf("[GetInstanceContainerLogs] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("[GetInstanceContainerLogs] Request succeded, but the response body could not be read. Status code: %d", res.StatusCode))
	}
	return string(body), nil
}

// GetInstanceSubmitLogsByInstanceID Method to retrieve instance submit logs
func (ac *Client) GetInstanceSubmitLogsByInstanceID(instanceID string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/instances/%s/submit-logs", ac.APIHost, instanceID), nil)
	if err != nil {
		return "", errors.Wrap(err, "[GetInstanceSubmitLogsByInstanceID] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "[GetInstanceSubmitLogsByInstanceID] Failed to execute GET request")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("[GetInstanceSubmitLogsByInstanceId] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return "", fmt.Errorf("[GetInstanceSubmitLogsByInstanceId] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("[GetInstanceSubmitLogsByInstanceId] Request succeded, but the response body could not be read. Status code: %d", res.StatusCode))
	}
	return string(body), nil
}
