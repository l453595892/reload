# ReloadController
目前方法是在部署reloadController以后可以主动更新deployment，statefulset，daemonset但是也有弊端，会无选择直接更新

## Template
使用Dockerfile编译打包好以后发布在k8s集群内部，或者更改kubeconfig路径，在本地编译测试

## TODO
根据标签选择相对应的资源进行更新

## 联系作者
WeChat:LWJ934110