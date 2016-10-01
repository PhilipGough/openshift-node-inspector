package main

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/philipgough/openshift-node-inspector/cmd"
	"github.com/philipgough/openshift-node-inspector/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strings"
)

const basePath string = "/tmp/oni/"


type DeploymentConfigPort struct {
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
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
		Run: func(cmnd *cobra.Command, args []string) {
			objects := []string{"svc", "dc"}
			for _, value := range objects {
				utils.ValidateInput(args[0], value)
				utils.SaveCleanFile(args[0], value, "/clean")
			}
			cmd.CreateDebugService(args[0], debugPort)
			//createDebugDc(args[0], debugPort)
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
