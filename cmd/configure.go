package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.build.ge.com/predix-data-services/predix-insights-go-sdk/predixinsights"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     "configure",
	Short:   "Login to Predix Insights",
	Long:    `Login and configure Predix Insights.`,
	Example: "  pi configure --interactive\n  pi configure --APIHost MY_API_HOST --TenantID MY_TENANT_ID --IssuerID MY_ISSUER_ID --ClientID MY_CLIENT_ID --ClientSecret MY_CLIENT_SECRET",
	Run: func(cmd *cobra.Command, args []string) {
		err := getMissingRequiredParams(loginPI)
		if err != nil {
			return
		}
		_, err = login()
		if err != nil {
			fmt.Println("authentication error err=" + err.Error())
			return
		}
		// maybe write entire viper in the future (will container all params)
		cleanup(loginPI)

		fmt.Println("login success")
	},
}

func login() (predixinsights.Client, error) {
	if loginPI.V.GetString("APIHost") == "" || loginPI.V.GetString("TenantID") == "" || loginPI.V.GetString("IssuerID") == "" || loginPI.V.GetString("ClientID") == "" || loginPI.V.GetString("ClientSecret") == "" {
		return predixinsights.Client{}, errors.New("please configure the Predix Insights CLI\n\n$ pi configure -i")
	}
	client := predixinsights.Client{APIHost: loginPI.V.GetString("APIHost"), TenantID: loginPI.V.GetString("TenantID"), IssuerID: loginPI.V.GetString("IssuerID"), ClientID: loginPI.V.GetString("ClientID"), ClientSecret: loginPI.V.GetString("ClientSecret"), Token: loginPI.V.GetString("Token")}

	// set verbose mode for pi-go-sdk
	client.Verbose = viper.GetBool("verbose")

	err := client.RefreshAuthToken()
	if err != nil {
		return predixinsights.Client{}, err
	}

	// save params
	loginPI.V.Set("Token", client.Token)

	return client, nil
}

func cleanup(pi pi) {
	// reset to default
	pi.V.Set("tail", false)
	pi.V.Set("verbose", false)
	pi.V.Set("interactive", false)
	pi.V.Set("containerlogsink", 1)

	// ensure dir exists
	if _, err := os.Stat(filepath.Dir(viper.ConfigFileUsed())); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(viper.ConfigFileUsed()), os.ModePerm)
	}
	// ensure the file exists
	if _, err := os.Stat(viper.ConfigFileUsed()); os.IsNotExist(err) {
		err = ioutil.WriteFile(viper.ConfigFileUsed(), []byte{}, os.FileMode(0644))
		if err != nil {
			fmt.Println("error creating initial viper config file err= " + err.Error())
		}
	}

	// set configuration file
	pi.V.SetConfigFile(viper.ConfigFileUsed())

	// write viper to config file
	err := pi.V.WriteConfig()
	if err != nil {
		fmt.Println("error saving viper config file err= " + err.Error())
	}
	if viper.GetBool("verbose") {
		fmt.Println("Saving config file:", pi.V.ConfigFileUsed())
	}
}

// json pretty format/print
func prettyprint(b []byte) {
	buff, err := jsonBeautify(b)
	if err != nil {
		fmt.Printf("erorr prettyprint %s", err.Error())
	}
	fmt.Printf("%s\n", buff.Bytes())
}

// JSON Beautify
func jsonBeautify(b []byte) (bytes.Buffer, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	if err != nil {
		return out, err
	}
	return out, nil
}
