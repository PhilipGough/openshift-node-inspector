package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ValidateInput(cliArg string, objectType string) {
	fmt.Println("Validating input...")

	if cliArg == "millicore" {
		fmt.Println("Theres always one...")
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
