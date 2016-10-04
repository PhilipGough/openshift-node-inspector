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
var niSrc string
var commitHash string

func CreateDebugService(nodeComponent string, debugPort int) {
	objectType = "svc"
	component = nodeComponent
	port = debugPort
	createDebugSvcFile()

}

func createDebugSvcFile() {
	defer utils.CreateDebugObj("update", objectType, component)
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

func CreateDebugDeploymentConfig(nodeComponent string, debugPort int, imageName string, src string, hash string) {
	objectType = "dc"
	component = nodeComponent
	port = debugPort
	image = imageName
	niSrc = src
	commitHash = hash
	createDebugDcFile()

}

func createDebugDcFile() {
	defer utils.CreateDebugObj("update", objectType, component)
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
	array = utils.UpdateContainerPorts(array, port)
	jsonParsed.S("spec", "template", "spec", "containers").Index(0).Set(array, "ports")

	// Make the component name available as an environment variable to the container
	envArray, ok := jsonParsed.S("spec", "template", "spec", "containers").Index(0).S("env").Data().([]interface{})
	if !ok {
		fmt.Printf("Error parsing existing %s %s JSON file. Value not an array \n", component, objectType)
		os.Exit(2)
	}
	envArray = utils.AddComponentEnvVar(envArray, component)
	jsonParsed.S("spec", "template", "spec", "containers").Index(0).Set(envArray, "env")

	// Add the "command" Array to overwrite the Dockerfile CMD definition
	jsonParsed.S("spec", "template", "spec", "containers").Index(0).Set(utils.UpdateDockerCMD(), "command")

	// Add the node inspector source as a volume
	_, exists := jsonParsed.S("spec", "template", "spec", "volumes").Data().([]interface{})
	if !exists {
		jsonParsed.S("spec", "template", "spec").ArrayOfSize(1, "volumes")
		jsonParsed.S("spec", "template", "spec", "volumes").SetIndex(utils.CreateNodeInspectorVolume(niSrc, commitHash), 0)
	} else {
		jsonParsed.ArrayAppend(utils.CreateNodeInspectorVolume(niSrc, commitHash), "spec", "template", "spec", "volumes")
	}

	//Mount volume inside container
	volumeMntArray, ok := jsonParsed.S("spec", "template", "spec", "containers").Index(0).S("volumeMounts").Data().([]interface{})
	if !ok {
		jsonParsed.S("spec", "template", "spec", "containers").Index(0).ArrayOfSize(1, "volumeMounts")
	}
	volumeMntArray = utils.MountContainerVolume(volumeMntArray)
	jsonParsed.S("spec", "template", "spec", "containers").Index(0).Set(volumeMntArray, "volumeMounts")

	// Adjust health checks
	jsonParsed.S("spec", "template", "spec", "containers").Index(0).SetP(1000, "livenessProbe.initialDelaySeconds")

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
	err := exec.Command("oc", "delete", objectType, component).Run()

	if err != nil {
		fmt.Println(err)
		fmt.Printf("Error deleting old  %s  for %s. Exiting ...", objectType, component)
		os.Exit(2)
	}

	fmt.Printf("Existing %s for %s removed \n", objectType, component)
}
