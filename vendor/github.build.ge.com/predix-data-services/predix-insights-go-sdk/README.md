# predix-insights-go-sdk
A Go SDK for the Predix Insights Service

# Installation 
`go get github.build.ge.com/predix-data-services/predix-insights-go-sdk`

Note: The command will fail and say "no Go files in...". This is okay. 
You should be able to navigate to the directory in your gopath now:
`cd $GOPATH/src/github.build.ge.com/predix-data-services`

# GoDocs
* [Link](https://github.build.ge.com/predix-data-services/predix-insights-go-sdk/tree/master/predixinsights)

# Example Usage
```
package main

import "github.build.ge.com/predix-data-services/predix-insights-go-sdk/predixinsights"

func main() {
	uiHost := "https://andromeda-ui-dev.core.predixdatafabric.com"
	tenantID := "my_tenant"
	issuerID := "https:123123123.uaa.predix.io/oauth/token"
	clientID := "user_name"
	clientSecret := "password"

	// Create new predixinsights client
	pic := predixinsights.NewClient(uiHost, tenantID, issuerID, clientID, clientSecret)

	// Login to UI
	err := pic.RefreshCookie()
	if err != nil {
		return err
	}

	// Get a Flow for id 12312391723981723
	flow, err := pic.GetFlow("12312391723981723")
	if err != nil {
		return err
	}
}

```
