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
		Long:  `Debug allows you to debug Node components using Node Inspector`,
		Run: func(cmnd *cobra.Command, args []string) {
			//objects := []string{"svc", "dc"}
			objects := []string{"dc"}
			for _, value := range objects {
				//utils.ValidateInput(args[0], value)
				utils.SaveCleanFile(args[0], value, "/clean")
			}
			cmd.CreateDebugDeploymentConfig(args[0], debugPort)
			//cmd.CreateDebugService(args[0], debugPort)

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
