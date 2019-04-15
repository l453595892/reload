package handler

import (
	"craftli.co/reload/pkg/util"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	v1beta1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (re *ResourceHandlerImpl) StatefulSetUpdate() {
	statefulSetList, err := re.getStatefulSetByResourceName()
	if err != nil {
		logrus.Errorf("Failed to list deployments %v", err)
	}
	for _, d := range statefulSetList.Items {
		envs := d.Spec.Template.Spec.Containers[0].Env
		if updateEnvVar(envs, util.UpdateUUID, xid.New().String()) == util.NotEnvFound {
			envs = append(envs, corev1.EnvVar{
				Name:  util.UpdateUUID,
				Value: xid.New().String(),
			})
		}
		d.Spec.Template.Spec.Containers[0].Env = envs
		_, err := util.KubeClient.AppsV1beta1().StatefulSets(re.Namespace).Update(&d)
		if err != nil {
			logrus.Errorf("Failed to update statefulSet %v", err)
		}
	}
}

func (re *ResourceHandlerImpl) getStatefulSetByResourceName() (*v1beta1.StatefulSetList, error) {
	statefulSetList, err := util.KubeClient.AppsV1beta1().StatefulSets(re.Namespace).List(meta_v1.ListOptions{})
	var statefulSetItems []v1beta1.StatefulSet
	if err != nil {
		return statefulSetList, err
	}
	for _, item := range statefulSetList.Items {
		if len(item.Spec.Template.Spec.Volumes) == 0 {
			continue
		}
		for _, volume := range item.Spec.Template.Spec.Volumes {
			if volume.ConfigMap != nil {
				if volume.ConfigMap.Name == re.ResourceName {
					item.Spec.UpdateStrategy.Type = v1beta1.RollingUpdateStatefulSetStrategyType
					statefulSetItems = append(statefulSetItems, item)
				}
			}
		}
	}
	statefulSetList.Items = statefulSetItems
	return statefulSetList, err
}
