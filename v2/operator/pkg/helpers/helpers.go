package helpers

import "fmt"

func GetDeploymentName(name string) string {
	return fmt.Sprintf("%s-operator-deploy", name)
}

func GetConfigMapName(name string) string {
	return fmt.Sprintf("%s-operator-cm", name)
}

func GetServiceName(name string) string {
	return fmt.Sprintf("%s-operator-svc", name)
}

type Operation string

const (
	CREATION_OPERATION Operation = "create"
	DELETION_OPERATION Operation = "delete"
)
