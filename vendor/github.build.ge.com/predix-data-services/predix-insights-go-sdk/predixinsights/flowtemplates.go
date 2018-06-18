package predixinsights

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// FlowTemplate struct represents flow template information
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
	BlobPath    string         `json:"blobPath"`
	SparkArgs   SparkArguments `json:"sparkArguments"`
	Flows       []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"flows"`
}

// EncapsulatedSparkArgs struct encapsulates spark arguments
type EncapsulatedSparkArgs struct {
	SparkArgs SparkArguments `json:"sparkArguments"`
}

// SparkArguments struct represents spark arguments
type SparkArguments struct {
	ApplicationArgs     []string          `json:"applicationArgs,omitempty"`
	ClassName           string            `json:"className,omitempty"`
	DriverCores         int               `json:"driverCores,omitempty"`
	DriverMemory        string            `json:"driverMemory,omitempty"`
	ExecutorMemory      string            `json:"executorMemory,omitempty"`
	NumExecutors        int               `json:"numExecutors,omitempty"`
	SparkListeners      []string          `json:"sparkListeners,omitempty"`
	LogConf             string            `json:"logConf,omitempty"`
	ExecutorEnv         map[string]string `json:"executorEnv,omitempty"`
	DriverJavaOptions   string            `json:"driverJavaOptions,omitempty"`
	ExecutorJavaOptions string            `json:"executorJavaOptions,omitempty"`
	Confs               map[string]string `json:"confs,omitempty"`
	SystemProps         map[string]string `json:"systemProps,omitempty"`
	FileName            string            `json:"fileName,omitempty"`
	FrameworkName       string            `json:"frameworkName,omitempty"`
}

// TagsArray []string represents list of tags
type TagsArray []string

// GetAllFlowsByTemplateIDResponse struct represents list of flows along iwth some meta information
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

// FlowResponse struct represents flow information
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

// SaveTagsForFlowTemplateResponse struct represents response obtained when tags are updated
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

// LaunchResponse struct represents response of flow launch
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

// FlowTemplatesResponseWithMetadata struct representing response for retrieving flow templates with meta data
type FlowTemplatesResponseWithMetadata struct {
	Content          []FlowTemplate `json:"content"`
	Last             bool           `json:"last"`
	TotalElements    int            `json:"totalElements"`
	TotalPages       int            `json:"totalPages"`
	First            bool           `json:"first"`
	Sort             interface{}    `json:"sort"`
	NumberOfElements int            `json:"numberOfElements"`
	Size             int            `json:"size"`
	Number           int            `json:"number"`
}

// ListConfigFiles struct representing the response from list config files
type ListConfigFiles []struct {
	FileName        string `json:"fileName"`
	FileSize        int    `json:"fileSize"`
	LastUpdatedTime int64  `json:"lastUpdatedTime"`
	Directory       bool   `json:"directory"`
}

// GetFlowByTemplateIDAndFlowID Method to retrieve flow by templateid and flowid
func (ac *Client) GetFlowByTemplateIDAndFlowID(templateID string, flowID string) (FlowResponse, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s%s%s", ac.APIHost, "/api/v1/flow-templates/", templateID, "/flows/", flowID), nil)
	if err != nil {
		return FlowResponse{}, errors.Wrap(err, "[GetFlowByTemplateIDAndFlowID] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return FlowResponse{}, errors.Wrap(err, "[GetFlowByTemplateIDAndFlowID] Failed to execute GET request")
	}

	defer res.Body.Close()
	ac.dumpResponse(res)

	var flowResponse FlowResponse
	err = json.NewDecoder(res.Body).Decode(&flowResponse)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			return FlowResponse{}, errors.Wrap(err, "[GetFlowByTemplateIDAndFlowID] Empty Body")
		case err != nil:
			return FlowResponse{}, errors.Wrap(err, fmt.Sprintf("[GetFlowByTemplateIdAndFlowId] Failed : req: %v", req))
		}
	}
	return flowResponse, nil

}

// GetAllFlowsByTemplateID Method to get all flows for a particular template by templateID
func (ac *Client) GetAllFlowsByTemplateID(templateID string) (GetAllFlowsByTemplateIDResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s%s", ac.APIHost, "/api/v1/flow-templates/", templateID, "/flows"), nil)
	if err != nil {
		return GetAllFlowsByTemplateIDResponse{}, errors.Wrap(err, "[GetAllFlowsByTemplateID] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return GetAllFlowsByTemplateIDResponse{}, errors.Wrap(err, "[GetAllFlowsByTemplateID] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	var getAllFlowsByTemplateIDResponse GetAllFlowsByTemplateIDResponse
	err = json.NewDecoder(res.Body).Decode(&getAllFlowsByTemplateIDResponse)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			return GetAllFlowsByTemplateIDResponse{}, errors.Wrap(err, "[GetAllFlowsByTemplateID] Empty Body")
		case err != nil:
			return GetAllFlowsByTemplateIDResponse{}, errors.Wrap(err, fmt.Sprintf("[GetAllFlowsByTemplateID] Failed to get flows: req: %v", req))
		}
	}
	return getAllFlowsByTemplateIDResponse, nil
}

// PostFlowTemplate Method to post new Template
func (ac *Client) PostFlowTemplate(flowTemplateName, templateFileName, templateFilePath, version, desc, flowType string) (FlowTemplate, error) {
	fields := []string{"metadata"}
	values := []string{fmt.Sprintf("{\"version\":\"%s\",\"user\":\"%s\",\"name\":\"%s\",\"description\":\"%s\",\"type\":\"%s\",\"tags\":[]}", version, ac.ClientID, flowTemplateName, desc, flowType)}

	// Load file to buffer
	buffer, contentType, err := newFileUploadBuffer(templateFileName, templateFilePath, fields, values)
	if err != nil {
		return FlowTemplate{}, errors.Wrap(err, "[PostFlowTemplate] Failed to create file upload buffer")
	}

	// Create new POST flow-template reqest
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", ac.APIHost, flowTemplateResource), &buffer)
	if err != nil {
		return FlowTemplate{}, errors.Wrap(err, "[PostFlowTemplate] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Set("Content-Type", contentType)
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return FlowTemplate{}, errors.Wrap(err, "[PostFlowTemplate] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	switch {
	case res.StatusCode == 409:
		return FlowTemplate{}, ErrResourceAlreadyExists
	case res.StatusCode != 201:
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return FlowTemplate{}, errors.Wrap(err, fmt.Sprintf("[PostFlowTemplate] Post Flow template failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return FlowTemplate{}, fmt.Errorf("[PostFlowTemplate] Upload request returned %d. Body: %s", res.StatusCode, string(body))
	}

	var flowTemplateResp FlowTemplate
	err = json.NewDecoder(res.Body).Decode(&flowTemplateResp)
	if err != nil {
		return FlowTemplate{}, errors.Wrap(err, fmt.Sprintf("[PostFlowTemplate] Failed to deecode flow template response. Status code: %d", res.StatusCode))
	}

	return flowTemplateResp, nil
}

// PostFlowTemplateUsingAnalyticFilePath Method to post flow template using existing analytic file by providing it's blobpath
func (ac *Client) PostFlowTemplateUsingAnalyticFilePath(version, user, flowTemplateName, blobPath, desc, flowType string) (FlowTemplate, error) {

	payload := strings.NewReader(fmt.Sprintf("{\n\t\"version\": \"%s\", \n\t\"user\": \"%s\", \n\t\"name\": \"%s\", \n\t\"blobPath\": \"%s\", \n\t\"description\": \"%s\", \n\t\"type\": \"%s\" , \n\t\"tags\":[]\n}", version, user, flowTemplateName, blobPath, desc, flowType))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", ac.APIHost, flowTemplateResource), payload)
	if err != nil {
		return FlowTemplate{}, errors.Wrap(err, "[PostFlowTemplateUsingAnalyticFilePath] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Set("Content-Type", "application/json")
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return FlowTemplate{}, errors.Wrap(err, "[PostFlowTemplateUsingAnalyticFilePath] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	switch {
	case res.StatusCode == 409:
		return FlowTemplate{}, ErrResourceAlreadyExists
	case res.StatusCode != 201:
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return FlowTemplate{}, errors.Wrap(err, fmt.Sprintf("[PostFlowTemplateUsingAnalyticFilePath] Post Flow template failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return FlowTemplate{}, fmt.Errorf("[PostFlowTemplateUsingAnalyticFilePath] Upload request returned %d. Body: %s", res.StatusCode, string(body))
	}

	var flowTemplateResp FlowTemplate
	err = json.NewDecoder(res.Body).Decode(&flowTemplateResp)
	if err != nil {
		return FlowTemplate{}, errors.Wrap(err, fmt.Sprintf("[PostFlowTemplateUsingAnalyticFilePath] Failed to deecode flow template response. Status code: %d", res.StatusCode))
	}

	return flowTemplateResp, nil
}

// LaunchFlow Method to launch flow by templateID and flowID
func (ac *Client) LaunchFlow(flowTemplateID, flowID string) (LaunchResponse, error) {

	// Create new LaunchFlow request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/flow-templates/%s/flows/%s/launch", ac.APIHost, flowTemplateID, flowID), nil)
	if err != nil {
		return LaunchResponse{}, errors.Wrap(err, "[LaunchFlow] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return LaunchResponse{}, errors.Wrap(err, "[LaunchFlow] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 202 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return LaunchResponse{}, errors.Wrap(err, fmt.Sprintf("[LaunchFlow] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return LaunchResponse{}, fmt.Errorf("[LaunchFlow] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	// Get app_id
	var launchResp LaunchResponse
	err = json.NewDecoder(res.Body).Decode(&launchResp)
	if err != nil {
		return LaunchResponse{}, errors.Wrap(err, fmt.Sprintf("[LaunchFlow] Failed to decode app response. Status code: %d", res.StatusCode))
	}

	return launchResp, nil
}

// DeleteFlow Method to delete flow by templateID and flowID
func (ac *Client) DeleteFlow(flowTemplateID, flowID string) error {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/flow-templates/%s/flows/%s", ac.APIHost, flowTemplateID, flowID), nil)
	if err != nil {
		return errors.Wrap(err, "[DeleteFlow] Failed to create DELETE request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[DeleteFlow] Failed to execute DELETE request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 204 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[DeleteFlow] Delete failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[DeleteFlow] Delete returned %d. Body: %s", res.StatusCode, string(body))
	}

	return nil
}

// PostFlow Method to post flow
func (ac *Client) PostFlow(flowName, flowTemplateID string) (Flow, error) {

	// New create flow request
	var fr FlowRequest
	fr.Name = flowName
	flowBytes, err := json.Marshal(fr)
	if err != nil {
		return Flow{}, errors.Wrap(err, fmt.Sprintf("[PostFlow] Failed to marshal new flow request"))
	}

	reqReader := bytes.NewReader(flowBytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/flow-templates/%s/flows", ac.APIHost, flowTemplateID), reqReader)
	if err != nil {
		return Flow{}, errors.Wrap(err, "[PostFlow] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute and handle request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Flow{}, errors.Wrap(err, fmt.Sprintf("[PostFlow] Client request to andromeda UI failed"))
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 201 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return Flow{}, errors.Wrap(err, fmt.Sprintf("[PostFlow] Post Flow failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return Flow{}, fmt.Errorf("[PostFlow] Upload request returned %d. Body: %s", res.StatusCode, string(body))
	}

	var flowResp Flow
	err = json.NewDecoder(res.Body).Decode(&flowResp)
	if err != nil {
		return Flow{}, errors.Wrap(err, fmt.Sprintf("[PostFlow] Failed to decode flow response. Status code: %d", res.StatusCode))
	}

	return flowResp, nil
}

// GetFlowTemplate Method to get flowTemplate by flowTemplateID
func (ac *Client) GetFlowTemplate(flowTemplateID string) (FlowTemplate, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", ac.APIHost, "/api/v1/flow-templates/", flowTemplateID), nil)
	if err != nil {
		return FlowTemplate{}, errors.Wrap(err, "[GetFlowTemplate] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return FlowTemplate{}, errors.Wrap(err, "[GetFlowTemplate] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	var ftr FlowTemplate
	err = json.NewDecoder(res.Body).Decode(&ftr)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			return FlowTemplate{}, errors.Wrap(err, "[GetFlowTemplate] Empty Body in response")
		case err != nil:
			return FlowTemplate{}, errors.Wrap(err, fmt.Sprintf("[GetFlowTemplate] Failed to get flow template: req: %v", req))
		}
	}
	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return FlowTemplate{}, fmt.Errorf("[GetFlowTemplate] Request failed, and the response body could not be read. Status code: %d", res.StatusCode)
		}
		return FlowTemplate{}, fmt.Errorf("[GetFlowTemplate] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	return ftr, nil
}

// GetAllFlowTemplatesByPage Method to retrieve all flowTemplates by page number
func (ac *Client) GetAllFlowTemplatesByPage(maxPages int) ([]FlowTemplate, error) {
	var allFlowTemplates []FlowTemplate

	for page := 0; page < maxPages; page++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s%s?page=%v", ac.APIHost, "/api/v1/flow-templates", page), nil)
		if err != nil {
			return []FlowTemplate{}, errors.Wrap(err, "[GetAllFlowTemplatesByPage] Failed to create GET request")
		}
		req.Header.Add("predix-zone-id", ac.TenantID)
		req.Header.Add("authorization", ac.Token)
		ac.dumpRequest(req)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return []FlowTemplate{}, errors.Wrap(err, "[GetAllFlowTemplatesByPage] Failed to execute GET request")
		}
		defer res.Body.Close()
		ac.dumpResponse(res)

		var fsr FlowTemplatesResponse
		err = json.NewDecoder(res.Body).Decode(&fsr)
		if err != nil {
			return []FlowTemplate{}, errors.Wrap(err, fmt.Sprintf("[GetAllFlowTemplatesByPage] Failed to decode response"))
		}

		newFlowTemplates := fsr.Content
		allFlowTemplates = append(allFlowTemplates, newFlowTemplates...)
		if len(newFlowTemplates) == 0 {
			return allFlowTemplates, nil
		}
	}

	return allFlowTemplates, nil
}

// GetAllFlowTemplates Method to retrieve all flow templates
func (ac *Client) GetAllFlowTemplates() (FlowTemplatesResponseWithMetadata, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", ac.APIHost, "/api/v1/flow-templates/"), nil)
	if err != nil {
		return FlowTemplatesResponseWithMetadata{}, errors.Wrap(err, "[GetAllFlowTemplates] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return FlowTemplatesResponseWithMetadata{}, errors.Wrap(err, "[GetAllFlowTemplates] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return FlowTemplatesResponseWithMetadata{}, fmt.Errorf("[GetAllFlowTemplates] Request failed, and the response body could not be read. Status code: %d", res.StatusCode)
		}
		return FlowTemplatesResponseWithMetadata{}, fmt.Errorf("[GetAllFlowTemplates] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	var ftr FlowTemplatesResponseWithMetadata
	err = json.NewDecoder(res.Body).Decode(&ftr)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			return FlowTemplatesResponseWithMetadata{}, errors.Wrap(err, "[GetAllFlowTemplates] Empty Body in response")
		case err != nil:
			return FlowTemplatesResponseWithMetadata{}, errors.Wrap(err, fmt.Sprintf("[GetAllFlowTemplates] Failed to get flow template: req: %v", req))
		}
	}
	return ftr, nil
}

// GetFlowTemplateByName Method to retrieve flow Templates by name
func (ac *Client) GetFlowTemplateByName(flowTemplateName string) (FlowTemplatesResponseWithMetadata, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", ac.APIHost, "/api/v1/flow-templates?name=", flowTemplateName), nil)
	if err != nil {
		return FlowTemplatesResponseWithMetadata{}, errors.Wrap(err, "[GetFlowTemplateByName] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return FlowTemplatesResponseWithMetadata{}, errors.Wrap(err, "[GetFlowTemplateByName] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return FlowTemplatesResponseWithMetadata{}, fmt.Errorf("[GetFlowTemplateByName] Request failed, and the response body could not be read. Status code: %d", res.StatusCode)
		}
		return FlowTemplatesResponseWithMetadata{}, fmt.Errorf("[GetFlowTemplateByName] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	var ftr FlowTemplatesResponseWithMetadata
	err = json.NewDecoder(res.Body).Decode(&ftr)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			return FlowTemplatesResponseWithMetadata{}, errors.Wrap(err, "[GetFlowTemplateByName] Empty Body in response")
		case err != nil:
			return FlowTemplatesResponseWithMetadata{}, errors.Wrap(err, fmt.Sprintf("[GetFlowTemplateByName] Failed to get flow template: req: %v", req))
		}
	}
	return ftr, nil
}

// DeleteFlowTemplate Method to delete flowTemplate by ID
func (ac *Client) DeleteFlowTemplate(flowTemplateID string) error {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/flow-templates/%s", ac.APIHost, flowTemplateID), nil)
	if err != nil {
		return errors.Wrap(err, "[DeleteFlowTemplate] Failed to create DELETE request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[DeleteFlowTemplate] Failed to execute DELETE request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 204 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[DeleteFlowTemplate] Delete failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[DeleteFlowTemplate] Delete returned %d. Body: %s", res.StatusCode, string(body))
	}

	return nil
}

// GetTagsByFlowTemplateID Method to get all tags for flowTemplate by flowTemplateID
func (ac *Client) GetTagsByFlowTemplateID(flowTemplateID string) (TagsArray, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s%s", ac.APIHost, "/api/v1/flow-templates/", flowTemplateID, "/tags"), nil)
	if err != nil {
		return TagsArray{}, errors.Wrap(err, "[GetTagsByFlowTemplateID] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TagsArray{}, errors.Wrap(err, "[GetTagsByFlowTemplateID] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	var gettagsbyflowtemplateidresponse TagsArray
	err = json.NewDecoder(res.Body).Decode(&gettagsbyflowtemplateidresponse)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			return TagsArray{}, errors.Wrap(err, "[GetTagsByFlowTemplateID] Empty body in response")
		case err != nil:
			return TagsArray{}, errors.Wrap(err, fmt.Sprintf("[GetTagsByFlowTemplateId] Failed : req: %v", req))
		}
	}
	return gettagsbyflowtemplateidresponse, nil
}

// SaveTagsForFlowTemplate Method to save tags for flowTemplate
func (ac *Client) SaveTagsForFlowTemplate(flowTemplateID string, tagsarray TagsArray) (SaveTagsForFlowTemplateResponse, error) {
	tagArrayBytes, err := json.Marshal(tagsarray)
	if err != nil {
		return SaveTagsForFlowTemplateResponse{}, errors.Wrap(err, "[SaveTagsForFlowTemplate] Failed to marshal new Tags Array")
	}

	payLoad := bytes.NewReader(tagArrayBytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/flow-templates/%s/tags", ac.APIHost, flowTemplateID), payLoad)
	if err != nil {
		return SaveTagsForFlowTemplateResponse{}, errors.Wrap(err, "[SaveTagsForFlowTemplate] Create new SaveTagsForFlowTemplate request failed")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute and handle request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return SaveTagsForFlowTemplateResponse{}, errors.Wrap(err, "[SaveTagsForFlowTemplate] Client request to SaveTagsForFlowTemplate failed")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	var saveTagsForFlowTemplateResponse SaveTagsForFlowTemplateResponse
	err = json.NewDecoder(res.Body).Decode(&saveTagsForFlowTemplateResponse)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			fmt.Println("", saveTagsForFlowTemplateResponse)

			return SaveTagsForFlowTemplateResponse{}, errors.Wrap(err, "[SaveTagsForFlowTemplate] Client request to SaveTagsForFlowTemplate failed")
		case err != nil:
			return SaveTagsForFlowTemplateResponse{}, errors.Wrap(err, fmt.Sprintf("[SaveTagsForFlowTemplate] Failed : %v. req: %v", err, req))
		}
	}
	return saveTagsForFlowTemplateResponse, nil

}

// GetTagsForFlowByFlowTemplateIDAndFlowID Method to get tags for a particular flow
func (ac *Client) GetTagsForFlowByFlowTemplateIDAndFlowID(flowTemplateID string, flowID string) (TagsArray, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s%s%s%s", ac.APIHost, "/api/v1/flow-templates/", flowTemplateID, "/flows/", flowID, "/tags"), nil)
	if err != nil {
		return TagsArray{}, errors.Wrap(err, "[GetTagsForFlowByFlowTemplateIDAndFlowID] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return TagsArray{}, errors.Wrap(err, "[GetTagsForFlowByFlowTemplateIDAndFlowID] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	var tagsArray TagsArray
	err = json.NewDecoder(res.Body).Decode(&tagsArray)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct

			return TagsArray{}, errors.Wrap(err, "[GetTagsForFlowByFlowTemplateIDAndFlowID] Empty response body")
		case err != nil:
			return TagsArray{}, errors.Wrap(err, fmt.Sprintf("[GetTagsForFlowByFlowTemplateIdAndFlowId] Failed : req: %v", req))
		}
	}
	return tagsArray, nil
}

// SaveTagsForFlow Method to save tags for a particular flow
func (ac *Client) SaveTagsForFlow(flowTemplateID string, flowID string, tagsarray TagsArray) (FlowResponse, error) {
	tagArrayBytes, err := json.Marshal(tagsarray)
	if err != nil {
		return FlowResponse{}, errors.Wrap(err, "[SaveTagsForFlow] Failed to marshal new Tags Array")
	}

	payLoad := bytes.NewReader(tagArrayBytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/flow-templates/%s/flows/%s/tags", ac.APIHost, flowTemplateID, flowID), payLoad)
	if err != nil {
		return FlowResponse{}, errors.Wrap(err, "[SaveTagsForFlow] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute and handle request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return FlowResponse{}, errors.Wrap(err, "[SaveTagsForFlow] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	var flowResponse FlowResponse
	err = json.NewDecoder(res.Body).Decode(&flowResponse)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			return FlowResponse{}, errors.Wrap(err, "[SaveTagsForFlow] Response body is empty")
		case err != nil:
			return FlowResponse{}, errors.Wrap(err, fmt.Sprintf("[SaveTagsForFlow] Failed : req: %v", req))
		}
	}
	return flowResponse, nil

}

// UpdateFlowTemplateByFlowTemplateIDUsingNewZip Method to update existing flowTemplate by uploading new zip
func (ac *Client) UpdateFlowTemplateByFlowTemplateIDUsingNewZip(flowTemplateID, flowTemplateName, templateFileName, templateFilePath, version, desc, flowType string) error {
	fields := []string{"metadata"}
	values := []string{fmt.Sprintf("{\"version\":\"%s\",\"user\":\"%s\",\"name\":\"%s\",\"description\":\"%s\",\"type\":\"%s\",\"tags\":[]}", version, ac.ClientID, flowTemplateName, desc, flowType)}

	// Load file to buffer
	buffer, contentType, err := newFileUploadBuffer(templateFileName, templateFilePath, fields, values)
	if err != nil {
		return errors.Wrap(err, "[UpdateFlowTemplateByFlowTemplateIDUsingNewZip] Failed to create file upload buffer")
	}

	// Create new POST flow-template reqest
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s/%s", ac.APIHost, flowTemplateResource, flowTemplateID), &buffer)
	if err != nil {
		return errors.Wrap(err, "[UpdateFlowTemplateByFlowTemplateIDUsingNewZip] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Set("Content-Type", contentType)
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[UpdateFlowTemplateByFlowTemplateIDUsingNewZip] Failed to execute POST request")
	}

	defer res.Body.Close()
	ac.dumpResponse(res)

	switch {
	case res.StatusCode == 404:
		return fmt.Errorf("[UpdateFlowTemplateByFlowTemplateIdUsingNewZip] Updation of Flow template failed, No template with %s id found. Status code: %d", flowTemplateID, res.StatusCode)

	case res.StatusCode == 401:
		return fmt.Errorf("[UpdateFlowTemplateByFlowTemplateIdUsingNewZip] Updation of Flow template failed. Unauthorized. Status code: %d", res.StatusCode)

	case res.StatusCode != 202:
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[UpdateFlowTemplateByFlowTemplateIdUsingNewZip] Updation of Flow template failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[UpdateFlowTemplateByFlowTemplateIdUsingNewZip] Updation of Flow template failed. Status %d. Body: %s", res.StatusCode, string(body))
	}
	// Successful Since this request returns no body just returning nil
	return nil
}

// UpdateFlowTemplateByFlowTemplateIDChangeSparkArguments Method to update existing flowTemplate by changing spark arguments
func (ac *Client) UpdateFlowTemplateByFlowTemplateIDChangeSparkArguments(flowTemplateID string, encapsulatedsparkargs EncapsulatedSparkArgs) error {
	encapsulatedsparkargsBytes, err := json.Marshal(encapsulatedsparkargs)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("[UpdateFlowTemplateByFlowTemplateIdChangeSparkArguments] Failed to marshal EncapsulatedSparkArgs"))
	}

	payLoad := bytes.NewReader(encapsulatedsparkargsBytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/flow-templates/%s", ac.APIHost, flowTemplateID), payLoad)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("[UpdateFlowTemplateByFlowTemplateIdChangeSparkArguments] Failed to create POST request"))
	}

	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute and handle request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("[UpdateFlowTemplateByFlowTemplateIdChangeSparkArguments] Failed to marshal EncapsulatedSparkArgs"))
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	switch {
	case res.StatusCode == 404:
		return fmt.Errorf("[UpdateFlowTemplateByFlowTemplateIdChangeSparkArguments] Updation of Flow template failed, No template with %s id found. Status code: %d", flowTemplateID, res.StatusCode)

	case res.StatusCode == 401:
		return fmt.Errorf("[UpdateFlowTemplateByFlowTemplateIdChangeSparkArguments] Updation of Flow template failed. Unauthorized. Status code: %d", res.StatusCode)

	case res.StatusCode != 202:
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[UpdateFlowTemplateByFlowTemplateIdChangeSparkArguments] Updation of Flow template failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[UpdateFlowTemplateByFlowTemplateIdChangeSparkArguments] Updation of Flow template failed. Status %d. Body: %s", res.StatusCode, string(body))
	}
	// Successful Since this request returns no body just returning nil
	return nil

}

// UpdateFlowChangeSparkArguments Method to update flow by changing spark arguments
func (ac *Client) UpdateFlowChangeSparkArguments(flowTemplateID string, flowID string, encapsulatedsparkargs EncapsulatedSparkArgs) error {
	encapsulatedsparkargsBytes, err := json.Marshal(encapsulatedsparkargs)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("[UpdateFlowChangeSparkArguments] Failed to marshal EncapsulatedSparkArgs"))
	}

	payLoad := bytes.NewReader(encapsulatedsparkargsBytes)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/flow-templates/%s/flows/%s", ac.APIHost, flowTemplateID, flowID), payLoad)
	if err != nil {
		return errors.Wrap(err, "[UpdateFlowChangeSparkArguments] Failed to create GET request")
	}

	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute and handle request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[UpdateFlowChangeSparkArguments] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	switch {
	case res.StatusCode == 404:
		return fmt.Errorf("[UpdateFlowChangeSparkArguments] Updation of Flow failed, No flow with %s id found. Status code: %d", flowID, res.StatusCode)

	case res.StatusCode == 401:
		return fmt.Errorf("[UpdateFlowChangeSparkArguments] Updation of Flow failed. Unauthorized. Status code: %d", res.StatusCode)

	case res.StatusCode != 202:
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[UpdateFlowChangeSparkArguments] Updation of Flow failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[UpdateFlowChangeSparkArguments] Updation of Flow failed. Status %d. Body: %s", res.StatusCode, string(body))
	}
	// Successful Since this request returns no body just returning nil
	return nil

}

// UpdateFlowByFlowTemplateIDAndFlowIDAddConfigFile Method to update flow by adding config files
func (ac *Client) UpdateFlowByFlowTemplateIDAndFlowIDAddConfigFile(flowTemplateID, flowID string, fileDetails []FileDetails) error {

	// Load file to buffer
	buffer, contentType, err := newFileUploadBufferMultipleFiles(fileDetails)
	if err != nil {
		return errors.Wrap(err, "[UpdateFlowByFlowTemplateIDAndFlowIDAddConfigFile] Failed to create file upload buffer")
	}

	// Create new reqest
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s/%s/flows/%s/config", ac.APIHost, flowTemplateResource, flowTemplateID, flowID), &buffer)
	if err != nil {
		return errors.Wrap(err, "[UpdateFlowByFlowTemplateIDAndFlowIDAddConfigFile] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Set("Content-Type", contentType)
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[UpdateFlowByFlowTemplateIDAndFlowIDAddConfigFile] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 201 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[UpdateFlowByFlowTemplateIdAndFlowIdAddConfigFile] Response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[UpdateFlowByFlowTemplateIdAndFlowIdAddConfigFile] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	//Successful no Body hence returning nil
	return nil
}

// UpdateFlowByFlowTemplateIDAndFlowIDDeleteConfigFile Method to delete config file from flow using flowTemplateID, flowID and fileName of config file
func (ac *Client) UpdateFlowByFlowTemplateIDAndFlowIDDeleteConfigFile(flowTemplateID, flowID, fileName string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s/flows/%s/config?file=%s", ac.APIHost, flowTemplateResource, flowTemplateID, flowID, fileName), nil)
	if err != nil {
		return errors.Wrap(err, "[UpdateFlowByFlowTemplateIDAndFlowIDDeleteConfigFile] Failed to create DELETE request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[UpdateFlowByFlowTemplateIDAndFlowIDDeleteConfigFile] Failed to execute DELETE request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[UpdateFlowByFlowTemplateIDAndFlowIDDeleteConfigFile] Delete failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[UpdateFlowByFlowTemplateIDAndFlowIDDeleteConfigFile] Delete returned %d. Body: %s", res.StatusCode, string(body))
	}

	return nil
}

// DownloadConfigFileByFlowTemplateIDAndFlowID Method to download config file using flowtemplateID and flowID
func (ac *Client) DownloadConfigFileByFlowTemplateIDAndFlowID(flowTemplateID, flowID, fileName string) ([]KeyValuePair, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/%s/flows/%s/config?file=%s", ac.APIHost, flowTemplateResource, flowTemplateID, flowID, fileName), nil)
	if err != nil {
		return []KeyValuePair{}, errors.Wrap(err, "[DownloadConfigFileByFlowTemplateIDAndFlowID] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)

	ac.dumpRequest(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []KeyValuePair{}, errors.Wrap(err, "[DownloadConfigFileByFlowTemplateIDAndFlowID] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)
	var tempresponse map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&tempresponse)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			return []KeyValuePair{}, errors.Wrap(err, fmt.Sprintf("[DownloadConfigFileByFlowTemplateIDAndFlowID] Response body is empty: req: %v", req))
		case err != nil:
			return []KeyValuePair{}, errors.Wrap(err, fmt.Sprintf("[DownloadConfigFileByFlowTemplateIDAndFlowID] Failed to get config file: req: %v", req))
		}
	}

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []KeyValuePair{}, fmt.Errorf("[DownloadConfigFileByFlowTemplateIDAndFlowID] Request failed, and the response body could not be read. Status code: %d", res.StatusCode)
		}
		return []KeyValuePair{}, fmt.Errorf("[DownloadConfigFileByFlowTemplateIDAndFlowID] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	response := []KeyValuePair{}
	for key, value := range tempresponse {
		response = append(response, KeyValuePair{key, value})
	}
	return response, nil
}

// ListConfigFileByFlowTemplateIDAndFlowID Method to retrieve list of all config files using flowTemplateID and flowID
func (ac *Client) ListConfigFileByFlowTemplateIDAndFlowID(flowTemplateID, flowID string) (ListConfigFiles, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s/%s/flows/%s/config", ac.APIHost, flowTemplateResource, flowTemplateID, flowID), nil)
	if err != nil {
		return ListConfigFiles{}, errors.Wrap(err, "[ListConfigFileByFlowTemplateIDAndFlowID] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)

	ac.dumpRequest(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ListConfigFiles{}, errors.Wrap(err, "[ListConfigFileByFlowTemplateIDAndFlowID] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)
	var listConfigFiles ListConfigFiles
	err = json.NewDecoder(res.Body).Decode(&listConfigFiles)
	if err != nil {
		switch {
		case err == io.EOF:
			// empty body, return empty struct
			return ListConfigFiles{}, errors.Wrap(err, fmt.Sprintf("[ListConfigFileByFlowTemplateIDAndFlowID] Response body is empty: req: %v", req))
		case err != nil:
			return ListConfigFiles{}, errors.Wrap(err, fmt.Sprintf("[ListConfigFileByFlowTemplateIDAndFlowID] Failed to list config files: req: %v", req))
		}
	}

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return ListConfigFiles{}, fmt.Errorf("[ListConfigFileByFlowTemplateIDAndFlowID] Request failed, and the response body could not be read. Status code: %d", res.StatusCode)
		}
		return ListConfigFiles{}, fmt.Errorf("[ListConfigFileByFlowTemplateIDAndFlowID] Request returned %d. Body: %s", res.StatusCode, string(body))
	}
	return listConfigFiles, nil
}
