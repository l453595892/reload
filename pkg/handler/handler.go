package handler

type ObjectHandler interface {
	Handler() error
}

type ResourceHandler interface {
	DeploymentUpdate()
}
