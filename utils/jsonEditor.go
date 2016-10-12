package utils

import (
	"strconv"
)

type VolumeDefinition struct {
	Name string              `json:"name"`
	Src  VolumeSrcDefinition `json:"gitRepo"`
}

type VolumeSrcDefinition struct {
	URL    string `json:"repository"`
	Commit string `json:"revision, omitempty"`
}

func UpdateContainerPorts(parsedArray []interface{}, port int) []interface{} {
	niPortMap := map[string]interface{}{"containerPort": port, "protocol": "TCP"}
	parsedArray = append(parsedArray, niPortMap)
	debugListenerPortMap := map[string]interface{}{"containerPort": 5858, "protocol": "TCP"}
	parsedArray = append(parsedArray, debugListenerPortMap)
	return parsedArray
}

func UpdateDockerCMD() []string {
	commands := []string{"bash", "-c", "bash /tmp/openshift-node-inspector-src/start.sh"}
	return commands
}

func AddEnvVars(parsedArray []interface{}, component string, port int) []interface{} {
	componentEnvVar := map[string]string{"name": "ONI_COMPONENT", "value": component}
	parsedArray = append(parsedArray, componentEnvVar)
	portEnvVar := map[string]string{"name": "ONI_DEBUG_PORT", "value": strconv.Itoa(port)}
	parsedArray = append(parsedArray, portEnvVar)
	return parsedArray
}

func CreateNodeInspectorVolume(src string, hash string) VolumeDefinition {
	volumeSrc := VolumeSrcDefinition{URL: src, Commit: hash}
	volumeSpec := VolumeDefinition{Name: "node-inspector-src", Src: volumeSrc}
	return volumeSpec
}

func MountContainerVolume(parsedArray []interface{}) []interface{} {
	nodeInspectorMount := map[string]string{"name": "node-inspector-src", "mountPath": "/tmp"}
	parsedArray = append(parsedArray, nodeInspectorMount)
	return parsedArray
}
