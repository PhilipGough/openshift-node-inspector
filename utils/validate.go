package utils

import (
  "os"
  "os/exec"
  "fmt"
  "strings"
)


func ValidateInput(cliArg string, objectType string) {
	objects := []string{"svc", "dc"}
	cmd := fmt.Sprintf("oc get %s | grep  -F  %s | awk '{print $1}'", objectType, cliArg)
	validate, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Print("Error validating input")
		os.Exit(1)
	}
	output := strings.TrimSpace(string(validate))
	if output != cliArg {
		fmt.Printf("No %s with the name  %s exists", objectType, cliArg)
		os.Exit(2)
	}

}
