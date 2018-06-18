package predixinsights

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// DependencyResponse struct representing Dependency information
type DependencyResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Tenant   string `json:"tenant"`
	Type     string `json:"type"`
	Deployed bool   `json:"deployed"`
}

// DependenciesResponse struct representing list of Dependency and some other meta information
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

// FileDetails struct representing File related information
type FileDetails struct {
	FileName     string
	FileLocation string
	Fields       []string
	Values       []string
}

// DependencyDetails struct representing Dependency related information
type DependencyDetails struct {
	Type         string
	FileName     string
	FileLocation string
}

// GetAllDependencies Method to retrieve all dependencies
func (ac *Client) GetAllDependencies() (DependenciesResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/dependencies/", ac.APIHost), nil)
	if err != nil {
		return DependenciesResponse{}, errors.Wrap(err, "[GetAllDependencies] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return DependenciesResponse{}, errors.Wrap(err, "[GetAllDependencies] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return DependenciesResponse{}, errors.Wrap(err, fmt.Sprintf("[GetAllDependencies] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return DependenciesResponse{}, fmt.Errorf("[GetAllDependencies] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	// Get app_id
	var dependenciesResponse DependenciesResponse
	err = json.NewDecoder(res.Body).Decode(&dependenciesResponse)
	if err != nil {
		return DependenciesResponse{}, errors.Wrap(err, fmt.Sprintf("[GetAllDependencies] Failed to decode response. Status code: %d", res.StatusCode))
	}
	return dependenciesResponse, nil

}

// GetDependencyByID Method to retrieve dependency by ID
func (ac *Client) GetDependencyByID(dependencyID string) (DependencyResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/dependencies/%s", ac.APIHost, dependencyID), nil)
	if err != nil {
		return DependencyResponse{}, errors.Wrap(err, "[GetDependencyByID] Failed to create GET request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return DependencyResponse{}, errors.Wrap(err, "[GetDependencyByID] Failed to execute GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return DependencyResponse{}, errors.Wrap(err, fmt.Sprintf("[GetDependencyByID] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return DependencyResponse{}, fmt.Errorf("[GetDependencyByID] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	// Get app_id
	var dependencyResponse DependencyResponse
	err = json.NewDecoder(res.Body).Decode(&dependencyResponse)
	if err != nil {
		return DependencyResponse{}, errors.Wrap(err, fmt.Sprintf("[GetDependencyByID] Failed to decode response. Status code: %d", res.StatusCode))
	}
	return dependencyResponse, nil
}

// PostDependency Method to post new dependency
func (ac *Client) PostDependency(dependencyType, dependencyFileName, dependencyFileLocation string) ([]DependencyResponse, error) {

	// Load file to buffer
	fields := []string{"metadata"}

	values := []string{fmt.Sprintf("{\"type\":\"%s\"}", dependencyType)}

	buffer, contentType, err := newFileUploadBuffer(dependencyFileName, dependencyFileLocation, fields, values)
	if err != nil {
		return []DependencyResponse{}, errors.Wrap(err, "[PostDependency] Failed to create File upload buffer")
	}

	// Create new POST Reqest
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/dependencies/", ac.APIHost), &buffer)
	if err != nil {
		return []DependencyResponse{}, errors.Wrap(err, "[PostDependency] Failed to create POST request")
	}

	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Set("Content-Type", contentType)
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []DependencyResponse{}, errors.Wrap(err, "[PostDependency] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	switch {
	case res.StatusCode == 409:
		return []DependencyResponse{}, fmt.Errorf("[PostDependency] Post Dependency failed. The Dependency with same name already exists. Status code: %d", res.StatusCode)

	case res.StatusCode != 201:
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []DependencyResponse{}, errors.Wrap(err, fmt.Sprintf("[PostDependency] Post Dependency failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return []DependencyResponse{}, fmt.Errorf("[PostDependency] Post Dependency request returned %d. Body: %s", res.StatusCode, string(body))
	}

	var dependencyResponse []DependencyResponse
	err = json.NewDecoder(res.Body).Decode(&dependencyResponse)
	if err != nil {
		return []DependencyResponse{}, errors.Wrap(err, fmt.Sprintf("[PostDependency] Failed to decode Post Dependency response. Status code: %d", res.StatusCode))
	}

	return dependencyResponse, nil
}

// PostMultipleDependencies Method to post multiple dependencies
func (ac *Client) PostMultipleDependencies(dependencies []DependencyDetails) ([]DependencyResponse, error) {

	var filesDetails []FileDetails
	for index, dependency := range dependencies {
		var filedetails FileDetails
		filedetails.Fields = []string{fmt.Sprintf("metadata%d", index)}
		filedetails.Values = []string{fmt.Sprintf("{\"type\":\"%s\"}", dependency.Type)}
		filedetails.FileName = dependency.FileName
		filedetails.FileLocation = dependency.FileLocation
		filesDetails = append(filesDetails, filedetails)
	}

	buffer, contentType, err := newFileUploadBufferMultipleFiles(filesDetails)
	if err != nil {
		return []DependencyResponse{}, errors.Wrap(err, "[PostMultipleDependencies] Failed to create files upload buffer")
	}

	// Create new POST Reqest
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/dependencies/", ac.APIHost), &buffer)
	if err != nil {
		return []DependencyResponse{}, errors.Wrap(err, "[PostMultipleDependencies] Failed to create POST request")
	}

	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Set("Content-Type", contentType)
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []DependencyResponse{}, errors.Wrap(err, "[PostMultipleDependencies] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)
	switch {
	case res.StatusCode == 409:
		return []DependencyResponse{}, fmt.Errorf("[PostMultipleDependencies] Post Dependency failed. The Dependency with same name already exists. Status code: %d", res.StatusCode)

	case res.StatusCode != 201:
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return []DependencyResponse{}, errors.Wrap(err, fmt.Sprintf("[PostMultipleDependencies] Post Dependency failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return []DependencyResponse{}, fmt.Errorf("[PostMultipleDependencies] Post Dependency request returned %d. Body: %s", res.StatusCode, string(body))
	}

	var dependencyResponse []DependencyResponse
	err = json.NewDecoder(res.Body).Decode(&dependencyResponse)
	if err != nil {
		return []DependencyResponse{}, errors.Wrap(err, fmt.Sprintf("[PostMultipleDependencies] Failed to decode Post Dependency response. Status code: %d", res.StatusCode))
	}

	return dependencyResponse, nil
}

// DeployDependencyByDependencyID Method to deploy dependency by ID
func (ac *Client) DeployDependencyByDependencyID(dependencyID string) error {
	// Create new POST Reqest
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/dependencies/deploy/%s", ac.APIHost, dependencyID), nil)
	if err != nil {
		return errors.Wrap(err, "[DeployDependencyByDependencyID] Failed to create POST request")
	}

	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Set("Content-Type", "application/json")
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[DeployDependencyByDependencyID] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[DeployDependencyByDependencyId] Deploy Dependency failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[DeployDependencyByDependencyId] Deploy Dependency request returned %d. Body: %s", res.StatusCode, string(body))
	}

	return nil
}

// DeployAllDependencies Method to deploy all dependencies
func (ac *Client) DeployAllDependencies() error {
	// Create new POST Reqest
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/dependencies/deploy/", ac.APIHost), nil)
	if err != nil {
		return errors.Wrap(err, "[DeployAllDependencies] Failed to create POST request")
	}

	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Set("Content-Type", "application/json")
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[DeployAllDependencies] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[DeployAllDependencies] Deploy Dependencies failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[DeployAllDependencies] Deploy Dependencies request returned %d. Body: %s", res.StatusCode, string(body))
	}

	return nil
}

// UnDeployAllDependencies Method to undeploy all dependencies
func (ac *Client) UnDeployAllDependencies() error {
	// Create new POST Reqest
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/dependencies/undeploy/", ac.APIHost), nil)
	if err != nil {
		return errors.Wrap(err, "[UnDeployAllDependencies] Failed to create POST request")
	}

	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Set("Content-Type", "application/json")
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[UnDeployAllDependencies] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 204 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[UnDeployAllDependencies] UnDeploy Dependencies failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[UnDeployAllDependencies] UnDeploy Dependencies request returned %d. Body: %s", res.StatusCode, string(body))
	}

	return nil
}

// UnDeployDependencyByDependencyID Method to undeploy dependency by ID
func (ac *Client) UnDeployDependencyByDependencyID(dependencyID string) error {
	// Create new POST Reqest
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/dependencies/undeploy/%s", ac.APIHost, dependencyID), nil)
	if err != nil {
		return errors.Wrap(err, "[UnDeployDependencyByDependencyID] Failed to create POST request")
	}

	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Set("Content-Type", "application/json")
	ac.dumpRequest(req)

	// Execute and handle requqest
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[UnDeployDependencyByDependencyID] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 204 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[UnDeployDependencyByDependencyId] Deploy Dependency failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[UnDeployDependencyByDependencyId] Deploy Dependency request returned %d. Body: %s", res.StatusCode, string(body))
	}

	return nil
}

// DeleteDependencyByID Method to delete dependency by ID
func (ac *Client) DeleteDependencyByID(dependencyID string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/dependencies/%s", ac.APIHost, dependencyID), nil)
	if err != nil {
		return errors.Wrap(err, "[DeleteDependencyByID] Failed to create DELETE request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[DeleteDependencyByID] Failed to execute DELETE request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 204 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[DeleteDependencyById] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[DeleteDependencyById] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	return nil
}
