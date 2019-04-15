package handler

import (
	"craftli.co/reload/pkg/util"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (re *ResourceHandlerImpl) DeploymentUpdate() {
	deploymentList, err := re.getDeploymentByResourceName()
	if err != nil {
		logrus.Errorf("Failed to list deployments %v", err)
	}
	for _, d := range deploymentList.Items {
		envs := d.Spec.Template.Spec.Containers[0].Env
		if updateEnvVar(envs, util.UpdateUUID, xid.New().String()) == util.NotEnvFound {
			envs = append(envs, v1.EnvVar{
				Name:  util.UpdateUUID,
				Value: xid.New().String(),
			})
		}
		d.Spec.Template.Spec.Containers[0].Env = envs
		_, err := util.KubeClient.ExtensionsV1beta1().Deployments(re.Namespace).Update(&d)
		if err != nil {
			logrus.Errorf("Failed to update deployments %v", err)
		}
	}
}

func (re *ResourceHandlerImpl) getDeploymentByResourceName() (*v1beta1.DeploymentList, error) {
	deploymentList, err := util.KubeClient.ExtensionsV1beta1().Deployments(re.Namespace).List(meta_v1.ListOptions{})
	var deploymentItems []v1beta1.Deployment
	if err != nil {
		return deploymentList, err
	}
	for _, item := range deploymentList.Items {
		if len(item.Spec.Template.Spec.Volumes) == 0 {
			continue
		}
		for _, volume := range item.Spec.Template.Spec.Volumes {
			if volume.ConfigMap != nil {
				if volume.ConfigMap.Name == re.ResourceName {
					deploymentItems = append(deploymentItems, item)
				}
			}
		}
	}
	deploymentList.Items = deploymentItems
	return deploymentList, err
}
