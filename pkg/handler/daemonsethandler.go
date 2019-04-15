package handler

import (
	"craftli.co/reload/pkg/util"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (re *ResourceHandlerImpl) DaemonSetUpdate() {
	daemonSetList, err := re.getDaemonSetByResourceName()
	if err != nil {
		logrus.Errorf("Failed to list deployments %v", err)
	}
	for _, d := range daemonSetList.Items {
		envs := d.Spec.Template.Spec.Containers[0].Env
		if updateEnvVar(envs, util.UpdateUUID, xid.New().String()) == util.NotEnvFound {
			envs = append(envs, v1.EnvVar{
				Name:  util.UpdateUUID,
				Value: xid.New().String(),
			})
		}
		d.Spec.Template.Spec.Containers[0].Env = envs
		_, err := util.KubeClient.ExtensionsV1beta1().DaemonSets(re.Namespace).Update(&d)
		if err != nil {
			logrus.Errorf("Failed to update daemonSet %v", err)
		}
	}
}

func (re *ResourceHandlerImpl) getDaemonSetByResourceName() (*v1beta1.DaemonSetList, error) {
	daemonSetList, err := util.KubeClient.ExtensionsV1beta1().DaemonSets(re.Namespace).List(meta_v1.ListOptions{})
	var daemonSetItems []v1beta1.DaemonSet
	if err != nil {
		return daemonSetList, err
	}
	for _, item := range daemonSetList.Items {
		if len(item.Spec.Template.Spec.Volumes) == 0 {
			continue
		}
		for _, volume := range item.Spec.Template.Spec.Volumes {
			if volume.ConfigMap != nil {
				if volume.ConfigMap.Name == re.ResourceName {
					daemonSetItems = append(daemonSetItems, item)
				}
			}
		}
	}
	daemonSetList.Items = daemonSetItems
	return daemonSetList, err
}
