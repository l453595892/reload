package resource

import (
	"craftli.co/reload/pkg/handler"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"reflect"
)

type UpdateHandler struct {
	OldObj interface{}
	NewObj interface{}
}

func (up UpdateHandler) Handler() error {
	if up.OldObj == nil || up.NewObj == nil {
		logrus.Errorf("The Resource is nil")
	}
	if ok, resource := up.checkObject(); ok { //rollingupdate
		resource.DeploymentUpdate()
		resource.DaemonSetUpdate()
		resource.StatefulSetUpdate()
	}
	return nil
}

func (up *UpdateHandler) checkObject() (bool, *handler.ResourceHandlerImpl) {
	if cp, ok := up.OldObj.(*v1.ConfigMap); ok {
		if !reflect.DeepEqual(cp.Data, up.NewObj.(*v1.ConfigMap).Data) {
			return true, &handler.ResourceHandlerImpl{
				ResourceName: up.OldObj.(*v1.ConfigMap).Name,
				Namespace:    up.OldObj.(*v1.ConfigMap).Namespace,
			}
		}
	} else if st, ok := up.OldObj.(*v1.Secret); ok {
		if !reflect.DeepEqual(st.Data, up.NewObj.(*v1.Secret).Data) {
			return true, &handler.ResourceHandlerImpl{
				ResourceName: up.OldObj.(*v1.Secret).Name,
				Namespace:    up.OldObj.(*v1.Secret).Namespace,
			}
		}
	}
	return false, &handler.ResourceHandlerImpl{}
}
