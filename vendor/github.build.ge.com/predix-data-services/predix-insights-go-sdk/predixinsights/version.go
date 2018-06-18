package predixinsights

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const versionResource = "/api/v1/version"

// CheckStatus Method to check status of predix insights
func (ac *Client) CheckVersion() (string, error) {

	req, err := http.NewRequest("Get", fmt.Sprintf("%s%s", ac.APIHost, versionResource), nil)
	if err != nil {
		return "", errors.Wrap(err, "[CheckVersion] Failed to create GET request")
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "[CheckVersion] Failed to successfully make GET request")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("[CheckVersion] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return "", fmt.Errorf("[CheckVersion] Request returned %d. Body: %s", res.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("[CheckVersion] Request succeded, but the response body could not be read. Status code: %d", res.StatusCode))
	}
	return string(body), nil
}
