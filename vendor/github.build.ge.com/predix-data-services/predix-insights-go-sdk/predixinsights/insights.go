package predixinsights

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/pkg/errors"
)

const (
	flowTemplateResource = "/api/v1/flow-templates"
	dagResource          = "/api/v1/dags"
	flowResource         = "/api/v1/flows"
)

// Client struct contains Client information
type Client struct {
	APIHost      string
	TenantID     string
	IssuerID     string
	ClientID     string
	ClientSecret string
	cookie       string
	cookieMux    sync.Mutex
	Token        string
	Verbose      bool
}

// ArgsRequest struct represents arguments for spark job
type ArgsRequest struct {
	Name           string                 `json:"name"`
	SparkArguments map[string]interface{} `json:"sparkArguments"`
}

// FlowTemplatesResponse struct represents list of FlowTemplates
type FlowTemplatesResponse struct {
	Content []FlowTemplate `json:"content"`
}

// LatestInstanceDetails struct contains summary of latest instance of flow
type LatestInstanceDetails struct {
	Summary struct {
		Status    string `json:"status"`
		StartTime int64  `json:"startTime"`
	} `json:"summary"`
}

// UAAResponse struct contains UAA Token information
type UAAResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	Jti         string `json:"jti"`
}

// ErrResourceAlreadyExists error to be thrown if resource already exists
var ErrResourceAlreadyExists = errors.New("Resource already exists")

// NewClient Method to retrieve new client object
func NewClient(APIHost, TenantID, IssuerID, ClientID, ClientSecret string) *Client {
	return &Client{
		APIHost:      APIHost,
		TenantID:     TenantID,
		IssuerID:     IssuerID,
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		cookie:       "",
		cookieMux:    sync.Mutex{},
	}
}

// PostArguments Method to post spark arguments
func (ac *Client) PostArguments(flowName string, flowTemplateID string, flowID string, sparkArgs map[string]interface{}) error {
	argsReq := ArgsRequest{
		Name:           flowName,
		SparkArguments: sparkArgs,
	}

	argsBytes, err := json.Marshal(argsReq)
	if err != nil {
		return errors.Wrap(err, "[PostArguments] Failed to marshal arguments into json")
	}
	argsReader := bytes.NewReader(argsBytes)

	// Create new save args request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/flow-templates/%s/flows/%s", ac.APIHost, flowTemplateID, flowID), argsReader)
	if err != nil {
		return errors.Wrap(err, "[PostArguments] Failed to create POST request")
	}
	req.Header.Add("predix-zone-id", ac.TenantID)
	req.Header.Add("authorization", ac.Token)
	req.Header.Add("content-type", "application/json")
	ac.dumpRequest(req)

	// Execute and handle request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[PostArguments] Failed to execute POST request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 202 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[PostArguments] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[PostArguments] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	return nil
}

// RefreshAuthToken Method to refresh UAA Token
func (ac *Client) RefreshAuthToken() error {
	url := fmt.Sprintf("%s%s", ac.IssuerID, "?grant_type=client_credentials")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Wrap(err, "[RefreshAuthToken] Failed to create a GET request")
	}
	req.SetBasicAuth(ac.ClientID, ac.ClientSecret)
	ac.dumpRequest(req)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return errors.Wrap(err, "[RefreshAuthToken] Failed to execute a GET request")
	}

	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[RefreshAuthToken] RefreshAuthToken failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[RefreshAuthToken] RefreshAuthToken request returned %d. Body: %s", res.StatusCode, string(body))
	}

	var uaaResponse UAAResponse
	err = json.NewDecoder(res.Body).Decode(&uaaResponse)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("[RefreshAuthToken] Failed to decode response. Status code: %d", res.StatusCode))

	}
	ac.Token = fmt.Sprintf("bearer %s", uaaResponse.AccessToken)
	return nil
}
