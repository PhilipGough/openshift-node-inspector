package utils

import (
	"encoding/json"
	"fmt"
)

type Route struct {
	Kind    string    `json:"kind"`
	Version string    `json:"apiVersion"`
	Meta    RouteMeta `json:"metadata"`
	Spec    RouteSpec `json:"spec"`
}

type RouteMeta struct {
	Name string `json:"name"`
}

type RouteSpec struct {
	Target     RouteTarget     `json:"to"`
	PortTarget RoutePortTarget `json:"port"`
}

type RouteTarget struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

type RoutePortTarget struct {
	Target string `json:"targetPort"`
}

func RouteConstructor(component string) {
	defer CreateDebugObj("route", component)
	targetPort := RoutePortTarget{Target: "node-inspector"}
	routeTarget := RouteTarget{Kind: "Service", Name: component}
	routeSpec := RouteSpec{Target: routeTarget, PortTarget: targetPort}
	metaName := component + "-node-inspector"
	routeMeta := RouteMeta{Name: metaName}
	baseRoute := Route{Kind: "Route", Version: "v1", Meta: routeMeta, Spec: routeSpec}

	route, err := json.Marshal(baseRoute)
	if err != nil {
		fmt.Println("Unable to automate route creation")
		return
	}

	WriteDebugFile(string(route), component, "route")

}
