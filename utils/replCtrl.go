package utils

import (
	"os/exec"
)

func ScalePods(component string) error {
	scaleErr := exec.Command("oc", "scale", "dc", component, "--replicas=0").Run()
	return scaleErr

}
