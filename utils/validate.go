package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var restrictedList =  []string{"millicore", "gitlab-shell", "mysql", "redis", "ups", "memcached"}

func ValidateInput(cliArg string, objectType string) {
	fmt.Println("Validating input...")

	if isRestricted(cliArg) {
		fmt.Printf("%s is not a Node component", cliArg)
		os.Exit(1)
	}

	cmd := fmt.Sprintf("oc get %s | grep  -F  %s | awk '{print $1}'", objectType, cliArg)
	validate, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		fmt.Println("Error validating input")
		os.Exit(1)
	}
	output := strings.TrimSpace(string(validate))
	if output != cliArg {
		fmt.Printf("No %s with the name  %s exists \n", objectType, cliArg)
		os.Exit(2)
	}

}


 func isRestricted(userInput string) bool {
 	for _, value := range restrictedList {
 		if value == userInput {
 			return true
 		}
 	}
 	return false
 }
