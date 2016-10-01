package cmd

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/philipgough/openshift-node-inspector/utils"
	"io/ioutil"
	"os"
	"os/exec"
)

type ServicePort struct {
	Name       string `json:"name"`
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
	TargetPort int    `json:"targetPort"`
}

const objectType string = "svc"

var component string
var port int

func CreateDebugService(nodeComponent string, debugPort int) {
	component = nodeComponent
	port = debugPort
	readCleanSvcFile()

}

func readCleanSvcFile() {
	defer deleteCleanObj()
	file, err := ioutil.ReadFile(utils.GetFilePath(component, objectType, "/clean"))
	if err != nil {
		fmt.Printf("Error reading %s file for %s. Exiting ...", objectType, component)
		os.Exit(2)
	}

	jsonParsed, err := gabs.ParseJSON(file)
	if err != nil {
		panic(err)
	}

	children, _ := jsonParsed.S("spec", "ports").Children()
	for _, child := range children {
		myMap := child.Data()
		dict, ok := myMap.(map[string]interface{})
		if ok {
			dict["name"] = component
		}
	}

	nodeInspectorPort := ServicePort{Name: "node-inspector", Protocol: "TCP", Port: port, TargetPort: port}
	jsonParsed.ArrayAppend(nodeInspectorPort, "spec", "ports")

	utils.WriteDebugFile(jsonParsed.String(), component, objectType)
}

func deleteCleanObj() {
	defer createDebugObj()
	err := exec.Command("oc", "delete", objectType, component).Run()

	if err != nil {
		fmt.Printf("Error deleting old  %s  for %s. Exiting ...", objectType, component)
		os.Exit(2)
	}

	fmt.Printf("Existing service for %s removed \n", component)
}

func createDebugObj() {
	fmt.Printf("Creating new debug %s for %s \n", objectType, component)
	path := utils.GetFilePath(component, objectType, "/debug")
	err := exec.Command("oc", "create", "-f", path).Run()
	if err != nil {
		fmt.Printf("Error creating new  %s  for %s. Exiting ...", objectType, component)
		os.Exit(2)
	}
	fmt.Printf("Debug %s for %s created \n", objectType, component)
}
