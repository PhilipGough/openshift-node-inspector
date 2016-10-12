package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

const basePath string = "/tmp/oni/"

func SaveCleanFile(component string, objectType string) {
	const prefix = "/clean"
	cmd := exec.Command("oc", "get", objectType, component, "-o", "json")
	os.MkdirAll(basePath+component, 0777)
	outfile, err := os.Create(basePath + component + prefix + objectType + ".json")

	if err != nil {
		fmt.Println("Error saving output to file. Exiting ...")
		os.Exit(2)
	}

	defer outfile.Close()
	cmd.Stdout = outfile

	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error saving file for %s %s. Exiting ...", objectType, component)
		os.Exit(2)
	}
	cmd.Wait()
	fmt.Printf("State saved for %s %s \n", objectType, component)
}

func GetFilePath(component string, objectType string, prefix string) string {
	return basePath + component + prefix + objectType + ".json"
}

func WriteDebugFile(contents string, component string, objectType string) {
	const prefix = "/debug"
	file, err := os.Create(basePath + component + prefix + objectType + ".json")

	if err != nil {
		fmt.Printf("Error creating debug file for %s %s. Aborting ...", objectType, component)
		os.Exit(1)
	}

	_, writeErr := io.WriteString(file, contents)

	if writeErr != nil {
		fmt.Printf("Error writing to debug file for %s %s. Aborting ...", objectType, component)
		os.Exit(1)
	}

	file.Close()
	fmt.Printf("Debug file created for %s %s \n", objectType, component)
}

func CreateDebugObj(command string, objectType string, component string) {
	fmt.Printf("Creating new debug %s for %s \n", objectType, component)
	path := GetFilePath(component, objectType, "/debug")
	err := exec.Command("oc", command, "-f", path).Run()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Error creating new  %s  for %s. \n", objectType, component)
	}
	fmt.Printf("Debug %s for %s created \n", objectType, component)
}
