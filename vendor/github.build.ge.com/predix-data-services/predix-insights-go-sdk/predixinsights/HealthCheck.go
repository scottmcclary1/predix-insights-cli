package predixinsights

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const statusResource = "/api/v1/status"

// CheckStatus Method to check status of predix insights
func (ac *Client) CheckStatus() error {

	req, err := http.NewRequest("Get", fmt.Sprintf("%s%s", ac.APIHost, statusResource), nil)
	if err != nil {
		return errors.Wrap(err, "[CheckStatus] Failed to create GET request")
	}
	ac.dumpRequest(req)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "[CheckStatus] Failed to successfully make GET request")
	}
	defer res.Body.Close()
	ac.dumpResponse(res)

	if res.StatusCode != 200 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("[CheckStatus] Request failed, and the response body could not be read. Status code: %d", res.StatusCode))
		}
		return fmt.Errorf("[CheckStatus] Request returned %d. Body: %s", res.StatusCode, string(body))
	}
	return nil
}
