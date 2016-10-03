package main

import (
	"fmt"
	"github.com/philipgough/openshift-node-inspector/cmd"
	"github.com/philipgough/openshift-node-inspector/utils"
	"github.com/spf13/cobra"
	"strings"
)

func main() {

	var debugPort int
	var image string

	var cmdDebug = &cobra.Command{
		Use:   "debug [component to debug]",
		Short: "Debug component with Node Inspector",
		Long:  `Debug allows you to debug Node components running on OPenShift using Node Inspector`,
		Run: func(cmnd *cobra.Command, args []string) {
			if len(args) > 0 {
				objects := []string{"svc", "dc"}
				for _, value := range objects {
					utils.ValidateInput(args[0], value)
					utils.SaveCleanFile(args[0], value)
				}
				cmd.CreateDebugService(args[0], debugPort)
				cmd.CreateDebugDeploymentConfig(args[0], debugPort, image)
			} else {
				fmt.Println("Component must be provided with debug command. \n Use --help for more info")
			}
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
	cmdDebug.Flags().StringVarP(&image, "image", "i", "", "Image to use (should include :tag) - Defaults to current deployment config")
	rootCmd.AddCommand(cmdDebug, cmdClean)
	rootCmd.Execute()

}
