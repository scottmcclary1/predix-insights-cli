package predixinsights

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// SubmitDetails struct represents submit command used and the hadoop environment information for a instance
type SubmitDetails struct {
	Command     string `json:"command"`
	Environment struct {
		HADOOPCONFDIR string `json:"HADOOP_CONF_DIR"`
		SPARKHOME     string `json:"SPARK_HOME"`
		SPARKCONFDIR  string `json:"SPARK_CONF_DIR"`
	} `json:"environment"`
}

// Summary struct represents flow instance summary information
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

// Flow struct represents Flow information
type Flow struct {
	ID                    string                `json:"id"`
	Created               int64                 `json:"created"`
	Updated               int64                 `json:"updated"`
	Version               string                `json:"version"`
	Name                  string                `json:"name"`
	Description           interface{}           `json:"description"`
	Type                  string                `json:"type"`
	Tags                  []interface{}         `json:"tags"`
	SparkArgs             SparkArguments        `json:"sparkArguments"`
	LatestInstanceDetails LatestInstanceDetails `json:"latestInstanceDetails"`
	FlowTemplate          FlowTemplate          `json:"flowTemplate"`
}

// FlowsResponse struct represents list of flow
type FlowsResponse struct {
	Content []Flow `json:"content"`
}

// FlowRequest struct contains name of flow used for making particular flow request
type FlowRequest struct {
	Name string `json:"name"`
}

// FlowDirectUploadResponse struct represents the response obtained when a flow is directly posted
type FlowDirectUploadResponse struct {
	ID           string      `json:"id"`
	Created      int64       `json:"created"`
	Updated      int64       `json:"updated"`
	Version      string      `json:"version"`
	User         string      `json:"user"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Type         string      `json:"type"`
	Tags         []string    `json:"tags"`
	FlowTemplate interface{} `json:"flowTemplate"`
}

// CreateFlowTemplateFromFlowResponse struct represents the response obtained when a flow template is created from an existing flow which was directly created
type CreateFlowTemplateFromFlowResponse struct {
	ID          string   `json:"id"`
	Created     int64    `json:"created"`
	Updated     int64    `json:"updated"`
	Version     string   `json:"version"`
	User        string   `json:"user"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Tags        []string `json:"tags"`
	BlobPath    string   `json:"blobPath"`
	Flows       []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"flows"`
}

// KeyValuePair struct representing key value pairs, useful when the structure of request or response is set of key value pairs with string key and no definite value (interface)
type KeyValuePair struct {
	Key   string
	Value interface{}
}

// GetAllFlows Method to retrieve flows by number of pages
func (ac *Client) GetAllFlows(maxPages int) ([]Flow, error) {
	var allFlows []Flow

	for page := 0; page < maxPages; page++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s%s?page=%v", ac.APIHost, "/api/v1/flows", page), nil)
		if err != nil {
			return []Flow{}, errors.Wrap(err, "[GetAllFlows] Failed to create GET request")
		}
		req.Header.Add("predix-zone-id", ac.TenantID)
		req.Header.Add("authorization", ac.Token)

		ac.dumpRequest(req)
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return []Flow{}, errors.Wrap(err, "[GetAllFlows] Failed to execute GET request")
		}
		defer res.Body.Close()
		ac.dumpResponse(res)

		var fsr FlowsResponse
		err = json.NewDecoder(res.Body).Decode(&fsr)
		if err != nil {
			return []Flow{}, errors.Wrap(err, fmt.Sprintf("[GetAllFlows] Failed to decode response"))
		}

		newFlows := fsr.Content
		allFlows = append(allFlows, newFlows...)
		if len(newFlows) == 0 {
			return allFlows, nil
		}
	}

	return allFlows, nil
}

// GetFlow Method to get flow by flowName
func (ac *Client) GetFlow(flowName string) (Flow, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", ac.APIHost, "/api/v1/flows/", flowName), nil)
	if err != nil {
		return Flow{}, errors.Wrap(err, "[GetFlow] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)

	ac.dumpRequest(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Flow{}, errors.Wrap(err, "[GetFlow] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	var f Flow
	err = json.NewDecoder(res.Body).Decode(&f)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			return Flow{}, errors.Wrap(err, fmt.Sprintf("[GetFlow] Response body is empty: req: %v", req))
		case err != nil:
			return Flow{}, errors.Wrap(err, fmt.Sprintf("[GetFlow] Failed to get flow template: req: %v", req))
		}
	}

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return Flow{}, fmt.Errorf("[GetFlow] Request failed, and the response body could not be read. Status code: %d", res.StatusCode)
		}
		return Flow{}, fmt.Errorf("[GetFlow] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	return f, nil
}

// StopFlow Method to stop flow by flowName
func (ac *Client) StopFlow(flowName string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/flows/%s/stop", ac.APIHost, flowName), nil)
	if err != nil {
		return errors.Wrap(err, "[StopFlow] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)

	ac.dumpRequest(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[StopFlow] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 202 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[StopFlow] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[StopFlow] Request returned %d. Body: %s", res.StatusCode, string(body))
	}
	return nil
}

// PostFlowDirectly Method to post flow directly without first uploading flowTemplate
func (ac *Client) PostFlowDirectly(flowName, flowFileName, flowFilePath, version, desc, flowType string) (FlowDirectUploadResponse, error) {
	fields := []string{"metadata"}
	values := []string{fmt.Sprintf("{\"version\":\"%s\",\"user\":\"%s\",\"name\":\"%s\",\"description\":\"%s\",\"type\":\"%s\",\"tags\":[]}", version, ac.ClientID, flowName, desc, flowType)}

	// Load file to buffer
	buffer, contentType, err := newFileUploadBuffer(flowFileName, flowFilePath, fields, values)
	if err != nil {
		return FlowDirectUploadResponse{}, errors.Wrap(err, "[PostFlowDirectly] Failed to create file upload buffer")
	}

	// Create new POST flow reqest
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", ac.APIHost, flowResource), &buffer)
	if err != nil {
		return FlowDirectUploadResponse{}, errors.Wrap(err, "[PostFlowDirectly] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Set("Content-Type", contentType)
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return FlowDirectUploadResponse{}, errors.Wrap(err, "[PostFlowDirectly] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 201 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return FlowDirectUploadResponse{}, errors.Wrap(err, fmt.Sprintf("[PostFlowDirectly] Post Flow Directly failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return FlowDirectUploadResponse{}, fmt.Errorf("[PostFlowDirectly] Upload request returned %d. Body: %s", res.StatusCode, string(body))
	}

	var flowDirectUploadResonse FlowDirectUploadResponse
	err = json.NewDecoder(res.Body).Decode(&flowDirectUploadResonse)
	if err != nil {
		return FlowDirectUploadResponse{}, errors.Wrap(err, fmt.Sprintf("[PostFlowDirectly] Failed to decode post flow directly response. Status code: %d", res.StatusCode))
	}

	return flowDirectUploadResonse, nil
}

// UpdateDirectFlowByFlowIDChangeAnalyticFile Method to update direct flow by changing analytic file
func (ac *Client) UpdateDirectFlowByFlowIDChangeAnalyticFile(flowID, description, flowFileName, flowFilePath string) (FlowDirectUploadResponse, error) {
	fields := []string{"metadata"}
	values := []string{fmt.Sprintf("{\"description\":\"%s\",\"tags\":[]}", description)}

	// Load file to buffer
	buffer, contentType, err := newFileUploadBuffer(flowFileName, flowFilePath, fields, values)
	if err != nil {
		return FlowDirectUploadResponse{}, errors.Wrap(err, "[UpdateDirectFlowByFlowIDChangeAnalyticFile] Failed to create file upload buffer")
	}

	// Create new POST reqest
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s/%s", ac.APIHost, flowResource, flowID), &buffer)
	if err != nil {
		return FlowDirectUploadResponse{}, errors.Wrap(err, "[UpdateDirectFlowByFlowIDChangeAnalyticFile] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Set("Content-Type", contentType)
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return FlowDirectUploadResponse{}, errors.Wrap(err, "[UpdateDirectFlowByFlowIDChangeAnalyticFile] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode < 200 || res.StatusCode > 300 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return FlowDirectUploadResponse{}, errors.Wrap(err, fmt.Sprintf("[UpdateDirectFlowByFlowIDChangeAnalyticFile] Post Flow Directly failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return FlowDirectUploadResponse{}, fmt.Errorf("[UpdateDirectFlowByFlowIDChangeAnalyticFile] Upload request returned %d. Body: %s", res.StatusCode, string(body))
	}

	var flowDirectUploadResonse FlowDirectUploadResponse
	err = json.NewDecoder(res.Body).Decode(&flowDirectUploadResonse)
	if err != nil {
		return FlowDirectUploadResponse{}, errors.Wrap(err, fmt.Sprintf("[UpdateDirectFlowByFlowIDChangeAnalyticFile] Failed to decode response. Status code: %d", res.StatusCode))
	}

	return flowDirectUploadResonse, nil

}

// CreateFlowTemplateFromFlow Method to create flow template from a directly uploaded flow
func (ac *Client) CreateFlowTemplateFromFlow(flowID string) (CreateFlowTemplateFromFlowResponse, error) {

	// Create new POST flow reqest
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s/%s/create-template", ac.APIHost, flowResource, flowID), nil)
	if err != nil {
		return CreateFlowTemplateFromFlowResponse{}, errors.Wrap(err, "[CreateFlowTemplateFromFlow] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return CreateFlowTemplateFromFlowResponse{}, errors.Wrap(err, "[CreateFlowTemplateFromFlow] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 201 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return CreateFlowTemplateFromFlowResponse{}, errors.Wrap(err, fmt.Sprintf("[CreateFlowTemplateFromFlow] Create Flow Template from directly uploaded Flow failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return CreateFlowTemplateFromFlowResponse{}, fmt.Errorf("[CreateFlowTemplateFromFlow] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	var createFlowTemplateFromFlowResponse CreateFlowTemplateFromFlowResponse
	err = json.NewDecoder(res.Body).Decode(&createFlowTemplateFromFlowResponse)
	if err != nil {
		return CreateFlowTemplateFromFlowResponse{}, errors.Wrap(err, fmt.Sprintf("[CreateFlowTemplateFromFlowResponse] Failed to decode response. Status code: %d", res.StatusCode))
	}

	return createFlowTemplateFromFlowResponse, nil
}

// DeleteFlowByFlowIDOnly Method to delte flow by flowID alone, useful if flow was created directly
func (ac *Client) DeleteFlowByFlowIDOnly(flowID string) error {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", ac.APIHost, flowResource, flowID), nil)
	if err != nil {
		return errors.Wrap(err, "[DeleteFlowByFlowIDOnly] Failed to create DELETE request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[DeleteFlowByFlowIDOnly] Failed to execute DELETE request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 204 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[DeleteFlowByFlowIDOnly] Delete failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[DeleteFlowByFlowIDOnly] Delete returned %d. Body: %s", res.StatusCode, string(body))
	}

	return nil
}

// UpdateFlowByFlowIDAddConfigFile Method to add config file to a flow using flowID
func (ac *Client) UpdateFlowByFlowIDAddConfigFile(flowID string, fileDetails []FileDetails) error {

	// Load file to buffer
	buffer, contentType, err := newFileUploadBufferMultipleFiles(fileDetails)
	if err != nil {
		return errors.Wrap(err, "[UpdateFlowByFlowIDAddConfigFile] Failed to create file upload buffer")
	}

	// Create new reqest
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s/%s/config", ac.APIHost, flowResource, flowID), &buffer)
	if err != nil {
		return errors.Wrap(err, "[UpdateFlowByFlowIDAddConfigFile] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Set("Content-Type", contentType)
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[UpdateFlowByFlowIDAddConfigFile] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[UpdateFlowByFlowIDAddConfigFile] Response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[UpdateFlowByFlowIDAddConfigFile] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	//Successful no Body hence returning nil
	return nil
}

// UpdateFlowByFlowIDDeleteConfigFile Method to delete config file from flow using flowID and fileName
func (ac *Client) UpdateFlowByFlowIDDeleteConfigFile(flowID, fileName string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s/config?file=%s", ac.APIHost, flowResource, flowID, fileName), nil)
	if err != nil {
		return errors.Wrap(err, "[UpdateFlowByFlowIDDeleteConfigFile] Failed to create DELETE request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[UpdateFlowByFlowIDDeleteConfigFile] Failed to execute DELETE request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[UpdateFlowByFlowIDDeleteConfigFile] Delete failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[UpdateFlowByFlowIDDeleteConfigFile] Delete returned %d. Body: %s", res.StatusCode, string(body))
	}
	return nil
}

// DownloadConfigFileByFlowID Method to download config file using flowID and fileName
func (ac *Client) DownloadConfigFileByFlowID(flowID, fileName string) ([]KeyValuePair, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/%s/config?file=%s", ac.APIHost, flowResource, flowID, fileName), nil)
	if err != nil {
		return []KeyValuePair{}, errors.Wrap(err, "[DownloadConfigFileByFlowID] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)

	ac.dumpRequest(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []KeyValuePair{}, errors.Wrap(err, "[DownloadConfigFileByFlowID] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)
	var tempresponse map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&tempresponse)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			return []KeyValuePair{}, errors.Wrap(err, fmt.Sprintf("[DownloadConfigFileByFlowID] Response body is empty: req: %v", req))
		case err != nil:
			return []KeyValuePair{}, errors.Wrap(err, fmt.Sprintf("[DownloadConfigFileByFlowID] Failed to get config file: req: %v", req))
		}
	}

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []KeyValuePair{}, fmt.Errorf("[DownloadConfigFileByFlowID] Request failed, and the response body could not be read. Status code: %d", res.StatusCode)
		}
		return []KeyValuePair{}, fmt.Errorf("[DownloadConfigFileByFlowID] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	response := []KeyValuePair{}
	for key, value := range tempresponse {
		response = append(response, KeyValuePair{key, value})
	}
	return response, nil
}

// ListConfigFilesByFlowID Method to retrieve list of all config files using flowID
func (ac *Client) ListConfigFilesByFlowID(flowID string) (ListConfigFiles, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/%s/config", ac.APIHost, flowResource, flowID), nil)
	if err != nil {
		return ListConfigFiles{}, errors.Wrap(err, "[ListConfigFilesByFlowID] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)

	ac.dumpRequest(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ListConfigFiles{}, errors.Wrap(err, "[ListConfigFilesByFlowID] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)
	var listConfigFiles ListConfigFiles
	err = json.NewDecoder(res.Body).Decode(&listConfigFiles)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			return ListConfigFiles{}, errors.Wrap(err, fmt.Sprintf("[ListConfigFilesByFlowID] Response body is empty: req: %v", req))
		case err != nil:
			return ListConfigFiles{}, errors.Wrap(err, fmt.Sprintf("[ListConfigFilesByFlowID] Failed to list config files: req: %v", req))
		}
	}

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return ListConfigFiles{}, fmt.Errorf("[ListConfigFilesByFlowID] Request failed, and the response body could not be read. Status code: %d", res.StatusCode)
		}
		return ListConfigFiles{}, fmt.Errorf("[ListConfigFilesByFlowID] Request returned %d. Body: %s", res.StatusCode, string(body))
	}
	return listConfigFiles, nil
}
