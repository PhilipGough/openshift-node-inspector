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

var objectType string
var component string
var port int
var image string

func CreateDebugService(nodeComponent string, debugPort int) {
	objectType = "svc"
	component = nodeComponent
	port = debugPort
	createDebugSvcFile()

}

func createDebugSvcFile() {
	defer deleteCleanObj()
	file, err := ioutil.ReadFile(utils.GetFilePath(component, objectType, "/clean"))
	if err != nil {
		fmt.Printf("Error reading %s file for %s. Exiting ...", objectType, component)
		os.Exit(2)
	}

	jsonParsed, err := gabs.ParseJSON(file)
	if err != nil {
		fmt.Printf("Error parsing existing %s %s JSON file. Exiting ... \n", component, objectType)
		os.Exit(2)
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

func CreateDebugDeploymentConfig(nodeComponent string, debugPort int, imageName string) {
	objectType = "dc"
	component = nodeComponent
	port = debugPort
	image = imageName
	createDebugDcFile()

}

func createDebugDcFile() {
	defer deleteCleanObj()

	file, err := ioutil.ReadFile(utils.GetFilePath(component, objectType, "/clean"))
	if err != nil {
		fmt.Printf("Error reading %s file for %s. Exiting ...", objectType, component)
		os.Exit(2)
	}

	jsonParsed, err := gabs.ParseJSON(file)
	if err != nil {
		fmt.Printf("Error parsing existing %s %s JSON file. Exiting ... \n", component, objectType)
		os.Exit(2)
	}

	// Add the additional definitions to the ports Array
	array, ok := jsonParsed.S("spec", "template", "spec", "containers").Index(0).S("ports").Data().([]interface{})

	if !ok {
		fmt.Printf("Error parsing existing %s %s JSON file. Value not an array \n", component, objectType)
		os.Exit(2)
	}
	niPortMap := map[string]interface{}{"containerPort": port, "protocol": "TCP"}
	array = append(array, niPortMap)
	debugListenerPortMap := map[string]interface{}{"containerPort": 5858, "protocol": "TCP"}
	array = append(array, debugListenerPortMap)
	jsonParsed.S("spec", "template", "spec", "containers").Index(0).Set(array, "ports")

	// Add the "command" Array to overwrite the Dockerfile CMD definition
	commands := []string{"bash", "-c", "chmod +x /tmp/openshift-node-inspector/start.sh && /tmp/openshift-node-inspector/start.sh"}
	jsonParsed.S("spec", "template", "spec", "containers").Index(0).Set(commands, "command")

	// Make the component name available as an environment variable to the container
	envArray, ok := jsonParsed.S("spec", "template", "spec", "containers").Index(0).S("env").Data().([]interface{})

	if !ok {
		fmt.Printf("Error parsing existing %s %s JSON file. Value not an array \n", component, objectType)
		os.Exit(2)
	}

	componentEnvVar := map[string]string{"name": "ONI_COMPONENT", "value": component}
	envArray = append(envArray, componentEnvVar)
	jsonParsed.S("spec", "template", "spec", "containers").Index(0).Set(envArray, "env")

	// Remove health checking to allow debugger to run without Pods reporting unreachable
	children, _ := jsonParsed.S("spec", "template", "spec", "containers").Index(0).ChildrenMap()
	for key, _ := range children {
		if key == "livenessProbe" || key == "readinessProbe" {
			delete(children, key)
		}
	}

	if image == "" {
		value, ok := jsonParsed.S("spec", "template", "spec", "containers").Index(0).S("image").Data().(string)
		if ok {
			jsonParsed.S("spec", "template", "spec", "containers").Index(0).Set(utils.ConfirmImage(value, true), "image")
		}
	} else {
		jsonParsed.S("spec", "template", "spec", "containers").Index(0).Set(utils.ConfirmImage(image, false), "image")
	}

	utils.WriteDebugFile(jsonParsed.String(), component, objectType)

}

func deleteCleanObj() {
	defer createDebugObj()
	err := exec.Command("oc", "delete", objectType, component).Run()

	if err != nil {
		fmt.Printf("Error deleting old  %s  for %s. Exiting ...", objectType, component)
		os.Exit(2)
	}

	fmt.Printf("Existing %s for %s removed \n", objectType, component)
}

func createDebugObj() {
	fmt.Printf("Creating new debug %s for %s \n", objectType, component)
	path := utils.GetFilePath(component, objectType, "/debug")
	err := exec.Command("oc", "create", "-f", path).Run()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Error creating new  %s  for %s. Exiting ...", objectType, component)
		os.Exit(2)
	}
	fmt.Printf("Debug %s for %s created \n", objectType, component)
}
