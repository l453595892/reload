package handler

import (
	"craftli.co/reload/pkg/util"
	"k8s.io/api/core/v1"
)

type ObjectHandler interface {
	Handler() error
}

type ResourceHandler interface {
	DeploymentUpdate()
	StatefulSetUpdate()
	DaemonSetUpdate()
}

type ResourceHandlerImpl struct {
	ResourceName string
	Namespace    string
}

func updateEnvVar(envs []v1.EnvVar, envar string, shaData string) int {
	for i := range envs {
		if envs[i].Name == envar {
			envs[i].Value = shaData
			return util.Updated
		}
	}
	return util.NotEnvFound
}
