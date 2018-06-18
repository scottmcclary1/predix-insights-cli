package main

import (
	"fmt"
	"os"

	"github.build.ge.com/predix-data-services/predix-insights-cli/cmd"
)

func main() {
	bashCompletion := os.Getenv("GENERATE_BASH_COMPLETION_FILE")
	if bashCompletion != "" {
		err := cmd.RootCmd.GenBashCompletionFile("scripts/pi_completion.sh")
		if err != nil {
			fmt.Println("error generating bash completion file err=" + err.Error())
		}
	}
	cmd.Execute()
}
