package main

type ServiceData struct {
	ServerName              string
	ServerType              string
	RegisterGrpcServerFunc  string
	RegisterHttpHandlerFunc string
}

type RegisterTemplateData struct {
	Services []*ServiceData
}
