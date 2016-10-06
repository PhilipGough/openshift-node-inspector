package cmd

import (
	"fmt"
	"github.com/philipgough/openshift-node-inspector/utils"
	"os"
	"os/exec"
)

func Cleanup(component string) {
	routeName := component + "-node-inspector"
	deleteObj("route", routeName)
	svcPath := utils.GetFilePath(component, "svc", "/clean")
	if _, err := os.Stat(svcPath); os.IsNotExist(err) {
		fmt.Println("No clean state exists to revert to. Exiting")
		os.Exit(2)
	}
	delSvcErr := deleteObj("svc", component)
	if delSvcErr == nil {
		exec.Command("oc", "create", "-f", svcPath).Run()
	}

	rollbackDc(component)

}

func deleteObj(objectType string, component string) error {
	err := exec.Command("oc", "delete", objectType, component).Run()
	if err != nil {
		fmt.Println(err)
		fmt.Printf("Error deleting old  %s  for %s", objectType, component)
	} else {
		fmt.Printf("Existing %s for %s removed \n", objectType, component)
	}
	return err
}

func rollbackDc(component string) {
	err := exec.Command("oc", "rollback", component).Run()
	if err != nil {
		fmt.Printf("Error rolling back deployment for %s", component)
	} else {
		fmt.Printf("Rolled back deployment for %s", component)
	}
}
