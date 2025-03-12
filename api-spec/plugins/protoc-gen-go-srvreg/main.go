package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"text/template"

	"github.com/golang/glog"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

//go:embed register.tmpl
var registerTmpl string

func main() {
	flag.Parse()
	defer glog.Flush()

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gp *protogen.Plugin) error {
		gp.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)

		gen := gp.NewGeneratedFile("service_registrer.pb.go", "")

		services := make([]*ServiceData, 0)
		for _, name := range gp.Request.FileToGenerate {
			file := gp.FilesByPath[name]
			pkg := file.GoImportPath

			for _, service := range file.Services {
				server := fmt.Sprintf("%sServer", service.GoName)
				services = append(services, &ServiceData{
					ServerName:              server,
					ServerType:              gen.QualifiedGoIdent(pkg.Ident(server)),
					RegisterGrpcServerFunc:  gen.QualifiedGoIdent(pkg.Ident(fmt.Sprintf("Register%s", server))),
					RegisterHttpHandlerFunc: gen.QualifiedGoIdent(pkg.Ident(fmt.Sprintf("Register%sHandler", service.GoName))),
				})
			}
		}

		data, err := readTemplate(registerTmpl, services)
		if err != nil {
			return err
		}
		gen.Write(data)

		return nil
	})
}

func readTemplate(path string, services []*ServiceData) ([]byte, error) {
	tmpl, err := template.New("splunk_enterprise").Parse(path)
	if err != nil {
		return nil, err
	}

	data := RegisterTemplateData{
		Services: services,
	}
	var buf bytes.Buffer

	err = tmpl.Execute(&buf, &data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
