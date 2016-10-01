package main

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"github.com/philipgough/openshift-node-inspector/utils"
)

const basePath string = "/tmp/oni/"

func deploydebug() {

}

type DeploymentConfigPort struct {
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
}

type ServicePort struct {
	Name       string `json:"name"`
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
	TargetPort int    `json:"targetPort"`
}

func createDebugSvc(component string, port int) {

	file, err := ioutil.ReadFile(basePath + component + "/cleansvc.json")
	if err != nil {
		panic(err)
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
	writeOutFile(jsonParsed.String(), component, "svc")

	deleteSvc, err := exec.Command("oc", "delete", "svc", component).Output()

	if err != nil {
		panic(err)
	}
	fmt.Printf(string(deleteSvc) + "  Creating new service ....")

	createSvc, err := exec.Command("oc", "create", "-f", basePath + component + "/debugsvc.json").Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf(string(createSvc))
}

func writeOutFile(contents string, component string, objectType string) {

	file, err := os.Create(basePath + component + "/debug" + objectType + ".json")

	if err != nil {
		panic(err)
	}

	n, err := io.WriteString(file, contents)

	if err != nil {
		fmt.Println(n)
		panic(err)
	}

	file.Close()
}

func saveClean(component string, objectType string) {
	cmd := exec.Command("oc", "get", objectType, component, "-o", "json")
	os.MkdirAll(basePath + component, 0777)
	outfile, err := os.Create(basePath + component + "/clean" + objectType + ".json")

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


func createDebugDc(component string, port int) {
	file, err := ioutil.ReadFile(basePath + component + "/cleandc.json")
	if err != nil {
		panic(err)
	}

	jsonParsed, err := gabs.ParseJSON(file)
	if err != nil {
		panic(err)
	}

	//spec templates spec containers
	containerPort := DeploymentConfigPort{Protocol: "TCP", Port: port}
	jsonParsed.ArrayAppend(containerPort, "spec", "template", "spec", "containers", "ports")
}

func main() {

	var debugPort int
	var image string

	var cmdDebug = &cobra.Command{
		Use:   "debug [component to debug]",
		Short: "Debug component with Node Inspector",
		Long:  `Debug allows you to debug Node components using Node Inspector`,
		Run: func(cmd *cobra.Command, args []string) {
			objects := []string{"svc", "dc"}
			for _, value := range objects {
				utils.ValidateInput(args[0], value)
				utils.SaveFile(args[0], value, "/clean")
			}
			createDebugSvc(args[0], debugPort)
			createDebugDc(args[0], debugPort)
		},
	}

	var cmdClean = &cobra.Command{
		Use:   "clean [string to print]",
		Short: "Revert to previous deployment configuration",
		Long: `print is for printing anything back to the screen.
    			For many years people have printed back to the screen.
    			`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Print: " + strings.Join(args, " "))
		},
	}

	var rootCmd = &cobra.Command{Use: "openshift-node-inspector"}
	cmdDebug.Flags().IntVarP(&debugPort, "port", "p", 9000, "Port to set debugger web host")
	cmdDebug.Flags().StringVarP(&image, "image", "i", "", "Image to use - Defaults to current")
	rootCmd.AddCommand(cmdDebug, cmdClean)
	rootCmd.Execute()

}
