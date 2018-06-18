package predixinsights

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// DAG Struct represents an airflow DAG
type DAG struct {
	ID       string `json:"id"`
	Created  int64  `json:"created"`
	Updated  int64  `json:"updated"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Deployed bool   `json:"deployed"`
}

// GetAllDAGsResponse Struct represents list of DAGs
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

// DAGResponse struct represents a response struct containing Airflow DAG information
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

type DAGStatuses struct {
	ScheduleInterval string   `json:"schedule_interval"`
	ActiveRuns       []string `json:"active_runs"`
	DagName          string   `json:"dag_name"`
	SuccessRuns      []string `json:"success_runs"`
	DagOwner         string   `json:"dag_owner"`
	DagID            string   `json:"dag_id"`
	FailedRuns       []string `json:"failed_runs"`
}

type SingleDAGStatus struct {
	DagName string        `json:"dag_name"`
	Dags    []DAGStatuses `json:"dags"`
}

type DAGRun struct {
	RunID           string `json:"run_id"`
	DagName         string `json:"dag_name"`
	DagTenantID     string `json:"dag_tenant_id"`
	EndDate         string `json:"end_date"`
	State           string `json:"state"`
	ExecutionDate   string `json:"execution_date"`
	ExternalTrigger string `json:"external_trigger"`
	DagOwner        string `json:"dag_owner"`
	DagID           string `json:"dag_id"`
	StartDate       string `json:"start_date"`
}

type SingleDAGRun struct {
	DagName         string `json:"dag_name"`
	DagTenantID     string `json:"dag_tenant_id"`
	Conf            string `json:"conf"`
	EndDate         string `json:"end_date"`
	State           string `json:"state"`
	ExecutionDate   string `json:"execution_date"`
	ExternalTrigger string `json:"external_trigger"`
	StartDate       string `json:"start_date"`
	DagID           string `json:"dag_id"`
}

type DagTaskInstance struct {
	TaskID        string `json:"task_id"`
	EndDate       string `json:"end_date"`
	RunID         string `json:"run_id"`
	ExecutionDate string `json:"execution_date"`
	DagRunState   string `json:"dag_run_state"`
	TaskState     string `json:"task_state"`
	StartDate     string `json:"start_date"`
}

type DagInfo struct {
	TaskInstances []DagTaskInstance `json:"task_instances"`
	DagTenantID   string            `json:"dag_tenant_id"`
	DagOwner      string            `json:"dag_owner"`
	DagID         string            `json:"dag_id"`
}

type AllTasks struct {
	DagName string               `json:"dag_name"`
	Dags    []map[string]DagInfo `json:"dags"`
}

type DagTaskInstanceByTaskID struct {
	JobID         string `json:"job_id"`
	EndDate       string `json:"end_date"`
	ExecutionDate string `json:"execution_date"`
	State         string `json:"state"`
	Duration      string `json:"duration"`
	StartDate     string `json:"start_date"`
	DagID         string `json:"dag_id"`
}

type DagInfoByTaskID struct {
	TaskInstances []DagTaskInstanceByTaskID `json:"task_instances"`
}

type TasksByTaskID struct {
	TaskID string                       `json:"task_id"`
	Dags   []map[string]DagInfoByTaskID `json:"dags"`
}

type TaskRunInfo struct {
	DagName       string `json:"dag_name"`
	TaskState     string `json:"task_state"`
	DagTenantID   string `json:"dag_tenant_id"`
	DagRunState   string `json:"dag_run_state"`
	TaskID        string `json:"task_id"`
	RunID         string `json:"run_id"`
	ExecutionDate string `json:"execution_date"`
	DagOwner      string `json:"dag_owner"`
	DagID         string `json:"dag_id"`
}

// GetAllDAGs Method to retrieve all DAGs
func (ac *Client) GetAllDAGs() (GetAllDAGsResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", ac.APIHost, "/api/v1/dags"), nil)
	if err != nil {
		return GetAllDAGsResponse{}, errors.Wrap(err, "[GetAllDAGs] Failed to create get request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return GetAllDAGsResponse{}, errors.Wrap(err, "[GetAllDAGs] Failed to execute get request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	var getAllDagsResponse GetAllDAGsResponse
	err = json.NewDecoder(res.Body).Decode(&getAllDagsResponse)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			return GetAllDAGsResponse{}, errors.Wrap(err, "[GetAllDAGs] Response body is empty")
		case err != nil:
			return GetAllDAGsResponse{}, errors.Wrap(err, "[GetAllDAGs] Failed")
		}
	}
	return getAllDagsResponse, nil
}

// PostDAG Method to post a new DAG
func (ac *Client) PostDAG(dagName, dagFileName, dagFilePath, version, desc, flowType string, dt DAGTemplate) (DAGResponse, error) {

	fields := []string{"metadata"}
	values := []string{fmt.Sprintf("{\"version\":\"%s\",\"user\":\"%s\",\"name\":\"%s\",\"description\":\"%s\",\"type\":\"%s\",\"tags\":[]}", version, ac.ClientID, dagName, desc, flowType)}

	// Load file to buffer
	buffer, contentType, err := newTemplatedUploadBuffer(dagFileName, dagFilePath, fields, values, dt)
	if err != nil {
		return DAGResponse{}, errors.Wrap(err, "[PostDAG] Failed to create file upload buffer")
	}

	// Create new POST flow-template reqest
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", ac.APIHost, dagResource), &buffer)
	if err != nil {
		return DAGResponse{}, errors.Wrap(err, "[PostDAG] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Set("Content-Type", contentType)
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return DAGResponse{}, errors.Wrap(err, "[PostDAG] Failed to execute POST requestr")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	switch {
	case res.StatusCode == 409:
		return DAGResponse{}, ErrResourceAlreadyExists
	case res.StatusCode != 201:
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return DAGResponse{}, errors.Wrap(err, fmt.Sprintf("[PostDAG] Post DAG template failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return DAGResponse{}, fmt.Errorf("[PostDAG] Upload request returned %d. Body: %s", res.StatusCode, string(body))
	}

	var dagResp DAGResponse
	err = json.NewDecoder(res.Body).Decode(&dagResp)
	if err != nil {
		return DAGResponse{}, errors.Wrap(err, fmt.Sprintf("[PostDAG] Failed to decode DAG response. Status code: %d", res.StatusCode))
	}

	return dagResp, nil

}

// UpdateDAG Method to update existing DAG
func (ac *Client) UpdateDAG(dagName, dagFileName, dagFilePath, version, desc, flowType string, dt DAGTemplate) error {
	fields := []string{"metadata"}
	values := []string{fmt.Sprintf("{\"version\":\"%s\",\"user\":\"%s\",\"name\":\"%s\",\"description\":\"%s\",\"type\":\"%s\",\"tags\":[]}", version, ac.ClientID, dagName, desc, flowType)}

	// Load file to buffer
	buffer, contentType, err := newTemplatedUploadBuffer(dagFileName, dagFilePath, fields, values, dt)
	if err != nil {
		return errors.Wrap(err, "[UpdateDAG] Failed to create upload buffer")
	}

	updateDagAPI := dagResource + "/" + dagName
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", ac.APIHost, updateDagAPI), &buffer)
	if err != nil {
		return errors.Wrap(err, "[UpdateDAG] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Set("Content-Type", contentType)
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[UpdateDAG] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 202 {
		return fmt.Errorf("[UpdateDAG] Failed. Status returned %d", res.StatusCode)
	}
	return nil
}

// DeleteDAG Method to delete DAG
func (ac *Client) DeleteDAG(name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s%s", ac.APIHost, "/api/v1/dags/", name), nil)
	if err != nil {
		return errors.Wrap(err, "[DeleteDAG] Failed to create DELETE request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[DeleteDAG] Failed to execute DELETE request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	switch responseStatus := res.StatusCode; responseStatus {
	case 401:
		body, err2 := ioutil.ReadAll(res.Body)
		if err2 != nil {
			return fmt.Errorf("[DeleteDAG] Unauthorized status code: %v. Unable to read ResponseBody", res.StatusCode)
		}
		return fmt.Errorf("[DeleteDAG] Unauthorized status code: %v. Response: %v", res.StatusCode, string(body))
	case 403:
		body, err2 := ioutil.ReadAll(res.Body)
		if err2 != nil {
			return fmt.Errorf("[DeleteDAG] Forbidden status code: %v. Unable to read ResponseBody", res.StatusCode)
		}
		return fmt.Errorf("[DeleteDAG] Forbidden status code: %v. Response: %v", res.StatusCode, string(body))
	case 204:
		return nil //successful
	case 200:
		return nil //successful
	case 202:
		return nil //successful
	default: //do nothing;
	}
	body, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		return fmt.Errorf("[DeleteDAG] Status code: %v. Unable to read ResponseBody", res.StatusCode)
	}
	return fmt.Errorf("[DeleteDAG] Status code: %v. Response: %v", res.StatusCode, string(body))
}

// GetDAG Method to retrieve DAG by name
func (ac *Client) GetDAG(name string) (DAGResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", ac.APIHost, "/api/v1/dags/", name), nil)
	if err != nil {
		return DAGResponse{}, errors.Wrap(err, "[GetDAG] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return DAGResponse{}, errors.Wrap(err, "[GetDAG] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return DAGResponse{}, errors.Wrap(err, fmt.Sprintf("[GetDAG] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return DAGResponse{}, fmt.Errorf("[GetDAG] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	var d DAGResponse
	err = json.NewDecoder(res.Body).Decode(&d)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			return DAGResponse{}, errors.Wrap(err, "[GetDAG] Response body is empty")
		case err != nil:
			return DAGResponse{}, errors.Wrap(err, fmt.Sprintf("[GetDAG] Failed to get flow template: req: %v res: %v", req, res))
		}
	}

	return d, nil
}

// DeployDAG Method to deploy DAG by name
func (ac *Client) DeployDAG(name string) error {
	fmt.Println()
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s/deploy", ac.APIHost, "/api/v1/dags/", name), nil)
	if err != nil {
		return errors.Wrap(err, "[DeployDAG] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[DeployDAG] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 202 {
		return fmt.Errorf("[DeployDAG] Request returned bad status code: %v. Response: %v", res.StatusCode, res)
	}

	return nil
}

// GetAllDAGsAllStatuses Method to get all dags statuses
func (ac *Client) GetAllDAGsAllStatuses() ([]DAGStatuses, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/statusall", ac.APIHost, dagResource), nil)
	if err != nil {
		return []DAGStatuses{}, errors.Wrap(err, "[GetAllDAGsAllStatuses] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []DAGStatuses{}, errors.Wrap(err, "[GetAllDAGsAllStatuses] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []DAGStatuses{}, fmt.Errorf("[GetAllDAGsAllStatuses] Request failed, and the response body could not be read. Status code: %d", res.StatusCode)
		}
		return []DAGStatuses{}, fmt.Errorf("[GetAllDAGsAllStatuses] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	bodyBytes, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		return []DAGStatuses{}, fmt.Errorf("[GetAllDAGsAllStatuses] Unable to read response body")
	}

	var data map[string]map[string]DAGStatuses
	err3 := json.Unmarshal(bodyBytes, &data)

	if err3 != nil {
		return []DAGStatuses{}, fmt.Errorf("[GetAllDAGsAllStatuses] Unable to unmarshal response body")
	}
	var resp []DAGStatuses
	for _, value := range data {
		for _, value2 := range value {
			var ds = DAGStatuses{}
			ds.ActiveRuns = value2.ActiveRuns
			ds.DagID = value2.DagID
			ds.DagName = value2.DagName
			ds.DagOwner = value2.DagOwner
			ds.FailedRuns = value2.FailedRuns
			ds.ScheduleInterval = value2.ScheduleInterval
			ds.SuccessRuns = value2.SuccessRuns
			resp = append(resp, ds)
		}
	}
	return resp, nil
}

func (ac *Client) GetDAGStatusByDAGName(dagName string) (SingleDAGStatus, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/status/%s", ac.APIHost, dagResource, dagName), nil)
	if err != nil {
		return SingleDAGStatus{}, errors.Wrap(err, "[GetDAGStatusByDAGName] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return SingleDAGStatus{}, errors.Wrap(err, "[GetDAGStatusByDAGName] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return SingleDAGStatus{}, fmt.Errorf("[GetDAGStatusByDAGName] Request failed, and the response body could not be read. Status code: %d", res.StatusCode)
		}
		return SingleDAGStatus{}, fmt.Errorf("[GetDAGStatusByDAGName] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	var singleDAGStatus SingleDAGStatus
	err = json.NewDecoder(res.Body).Decode(&singleDAGStatus)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			return SingleDAGStatus{}, errors.Wrap(err, fmt.Sprintf("[GetDAGStatusByDAGName] Response body is empty: req: %v", req))
		case err != nil:
			return SingleDAGStatus{}, errors.Wrap(err, fmt.Sprintf("[GetDAGStatusByDAGName] Failed to get flow template: req: %v", req))
		}
	}

	return singleDAGStatus, nil
}

func (ac *Client) GetRunsByDAGName(dagName string) ([]DAGRun, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/status/%s/runs", ac.APIHost, dagResource, dagName), nil)
	if err != nil {
		return []DAGRun{}, errors.Wrap(err, "[GetRunsByDAGName] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []DAGRun{}, errors.Wrap(err, "[GetRunsByDAGName] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []DAGRun{}, fmt.Errorf("[GetRunsByDAGName] Request failed, and the response body could not be read. Status code: %d", res.StatusCode)
		}
		return []DAGRun{}, fmt.Errorf("[GetRunsByDAGName] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	bodyBytes, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		return []DAGRun{}, fmt.Errorf("[GetRunsByDAGName] Unable to read response body")
	}

	var data map[string][]DAGRun
	err3 := json.Unmarshal(bodyBytes, &data)

	if err3 != nil {
		return []DAGRun{}, fmt.Errorf("[GetRunsByDAGName] Unable to unmarshal response body")
	}

	var resp []DAGRun
	for _, value := range data {
		for _, element := range value {
			var dr = DAGRun{}
			dr.DagID = element.DagID
			dr.DagName = element.DagName
			dr.DagOwner = element.DagOwner
			dr.DagTenantID = element.DagTenantID
			dr.EndDate = element.EndDate
			dr.ExecutionDate = element.ExecutionDate
			dr.ExternalTrigger = element.ExternalTrigger
			dr.RunID = element.RunID
			dr.StartDate = element.StartDate
			dr.State = element.State
			resp = append(resp, dr)
		}
	}
	return resp, nil
}

func (ac *Client) GetRunByDAGNameAndRunID(dagName, runID string) (SingleDAGRun, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/status/%s/runs/%s", ac.APIHost, dagResource, dagName, runID), nil)
	if err != nil {
		return SingleDAGRun{}, errors.Wrap(err, "[GetRunByDAGNameAndRunID] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return SingleDAGRun{}, errors.Wrap(err, "[GetRunByDAGNameAndRunID] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return SingleDAGRun{}, fmt.Errorf("[GetRunByDAGNameAndRunID] Request failed, and the response body could not be read. Status code: %d", res.StatusCode)
		}
		return SingleDAGRun{}, fmt.Errorf("[GetRunByDAGNameAndRunID] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	bodyBytes, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		return SingleDAGRun{}, fmt.Errorf("[GetRunByDAGNameAndRunID] Unable to read response body")
	}

	var data map[string]SingleDAGRun
	err3 := json.Unmarshal(bodyBytes, &data)

	if err3 != nil {
		return SingleDAGRun{}, fmt.Errorf("[GetRunByDAGNameAndRunID] Unable to unmarshal response body")
	}

	var resp SingleDAGRun
	for _, value := range data {
		resp.DagID = value.DagID
		resp.DagName = value.DagName
		resp.DagTenantID = value.DagTenantID
		resp.EndDate = value.EndDate
		resp.ExecutionDate = value.ExecutionDate
		resp.ExternalTrigger = value.ExternalTrigger
		resp.StartDate = value.StartDate
		resp.State = value.State
	}
	return resp, nil
}

func (ac *Client) GetAllTasksByDagName(dagName string) (AllTasks, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/status/%s/tasks", ac.APIHost, dagResource, dagName), nil)
	if err != nil {
		return AllTasks{}, errors.Wrap(err, "[GetAllTasksByDagName] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return AllTasks{}, errors.Wrap(err, "[GetAllTasksByDagName] Failed to execute GET request")
	}

	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return AllTasks{}, errors.Wrap(err, fmt.Sprintf("[GetAllTasksByDagName] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return AllTasks{}, fmt.Errorf("[GetAllTasksByDagName] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	// Get app_id
	var tasks AllTasks
	err = json.NewDecoder(res.Body).Decode(&tasks)
	if err != nil {
		return AllTasks{}, errors.Wrap(err, fmt.Sprintf("[GetAllTasksByDagName] Failed to decode response. Status code: %d", res.StatusCode))
	}

	return tasks, nil
}

func (ac *Client) GetAllTasksByDagNameAndTaskID(dagName, taskID string) (TasksByTaskID, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/status/%s/tasks/%s", ac.APIHost, dagResource, dagName, taskID), nil)
	if err != nil {
		return TasksByTaskID{}, errors.Wrap(err, "[GetAllTasksByDagNameAndTaskID] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TasksByTaskID{}, errors.Wrap(err, "[GetAllTasksByDagNameAndTaskID] Failed to execute GET request")
	}

	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return TasksByTaskID{}, errors.Wrap(err, fmt.Sprintf("[GetAllTasksByDagNameAndTaskID] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return TasksByTaskID{}, fmt.Errorf("[GetAllTasksByDagNameAndTaskID] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	// Get app_id
	var tasks TasksByTaskID
	err = json.NewDecoder(res.Body).Decode(&tasks)
	if err != nil {
		return TasksByTaskID{}, errors.Wrap(err, fmt.Sprintf("[GetAllTasksByDagNameAndTaskID] Failed to decode response. Status code: %d", res.StatusCode))
	}

	return tasks, nil
}

func (ac *Client) GetTaskRunInfo(dagName, taskID, runID string) (TaskRunInfo, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/status/%s/tasks/%s/runs/%s", ac.APIHost, dagResource, dagName, taskID, runID), nil)
	if err != nil {
		return TaskRunInfo{}, errors.Wrap(err, "[GetTaskRunInfo] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TaskRunInfo{}, errors.Wrap(err, "[GetTaskRunInfo] Failed to execute GET request")
	}

	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return TaskRunInfo{}, errors.Wrap(err, fmt.Sprintf("[GetTaskRunInfo] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return TaskRunInfo{}, fmt.Errorf("[GetTaskRunInfo] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	// Get app_id
	var taskRunInfo TaskRunInfo
	err = json.NewDecoder(res.Body).Decode(&taskRunInfo)
	if err != nil {
		return TaskRunInfo{}, errors.Wrap(err, fmt.Sprintf("[GetTaskRunInfo] Failed to decode response. Status code: %d", res.StatusCode))
	}

	return taskRunInfo, nil
}
