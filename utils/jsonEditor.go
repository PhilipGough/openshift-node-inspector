package utils

func UpdateContainerPorts(parsedArray []interface{}, port int) []interface{} {
	niPortMap := map[string]interface{}{"containerPort": port, "protocol": "TCP"}
	parsedArray = append(parsedArray, niPortMap)
	debugListenerPortMap := map[string]interface{}{"containerPort": 5858, "protocol": "TCP"}
	parsedArray = append(parsedArray, debugListenerPortMap)
	return parsedArray
}

func UpdateDockerCMD() []string {
	commands := []string{"bash", "-c", "chmod +x /tmp/openshift-node-inspector/start.sh && /tmp/openshift-node-inspector/start.sh"}
	return commands
}

func AddComponentEnvVar(parsedArray []interface{}, component string) []interface{} {
	componentEnvVar := map[string]string{"name": "ONI_COMPONENT", "value": component}
	parsedArray = append(parsedArray, componentEnvVar)
	return parsedArray
}
