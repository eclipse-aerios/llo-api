package config

import (
	"strings"
)

const SERVICE_NAME = "aeriOS LLO API"
const API_VERSION = "1.2.0"
const HEALTHY_STATUS = "HEALTHY"
const UNHEALTHY_STATUS = "UNHEALTHY"

const GVR_GROUP = "llo.aeros-project.eu"
const GVR_VERSION = "v1alpha1"

const CONTAINERD_CR = "ServiceComponentContainerd"
const DOCKER_CR = "ServiceComponentK8s"
const K8S_CR = "ServiceComponentDocker"
const DEFAULT_CR = K8S_CR
const DEFAULT_LLO_TYPE = "k8s"

// TODO namespace as path param
var Namespace string = "default"
var Status string = HEALTHY_STATUS

// TODO precheck if CRDs exists in the cluster

func GetSupportedCRs() (supportedCRs []string) {
	supportedCRs = []string{
		DOCKER_CR,
		K8S_CR,
	}
	return
}

// TODO return error if ""
func GetCR(name string) (cr string) {
	name = strings.ToLower(name)
	switch name {
	case "docker", "servicecomponentdocker":
		cr = "servicecomponentdockers"
	case "k8s", "servicecomponentk8s":
		cr = "servicecomponentk8s"
	// case "containerd", "servicecomponentcontainerd":
	// 	cr = "servicecomponentcontainerds"
	default:
		cr = ""
	}
	return
}
