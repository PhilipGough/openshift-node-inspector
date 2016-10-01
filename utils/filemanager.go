package utils

import (
	"os"
	"os/exec"
)

const basePath string = "/tmp/oni/"

// Clean or dirty is prefix
func SaveFile(component string, objectType string, prefix string) {
	cmd := exec.Command("oc", "get", objectType, component, "-o", "json")
	os.MkdirAll(basePath + component, 0777)
	outfile, err := os.Create(basePath + component + prefix + objectType + ".json")

	if err != nil {
		panic(err)
	}

	defer outfile.Close()
	cmd.Stdout = outfile

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	cmd.Wait()
}
