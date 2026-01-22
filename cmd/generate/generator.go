package main

import (
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

type Generator struct {
	file *descriptorpb.FileDescriptorProto
}

func (g *Generator) IsHaveService() bool {
	return len(g.file.GetService()) != 0
}

func (g *Generator) Generate() *pluginpb.CodeGeneratorResponse_File {
	file := g.file
	pkg := sanitizePackage(file.GetPackage())

	filename := file.GetName()
	outname := strings.TrimSuffix(filename, ".proto") + "_dispatcher.pb.go"
	outnames := strings.Split(outname, "/")

	outnames = append(outnames[:len(outnames)-1], pkg, outnames[len(outnames)-1])
	outname = strings.Join(outnames, "/")

	var sb strings.Builder

	sb.WriteString(g.Import())
	sb.WriteString(g.GenerateServices())

	return &pluginpb.CodeGeneratorResponse_File{
		Name:    proto.String(outname),
		Content: proto.String(sb.String()),
	}
}

func (g *Generator) Import() string {
	file := g.file
	var sb strings.Builder

	pkg := sanitizePackage(file.GetPackage())

	// import package
	gopkg := file.Options.GetGoPackage()
	gopkgs := strings.Split(gopkg, ";")
	gopkg = gopkgs[0]
	gopkgs = strings.Split(gopkg, "/")

	pkgImport := gopkgs[len(gopkgs)-1] + " " + `"` + gopkg + `"`

	sb.WriteString("package " + pkg + "\n\n")

	sb.WriteString("import (\n")
	sb.WriteString(`    "context"` + "\n")

	sb.WriteString(`    connect "connectrpc.com/connect"` + "\n")
	sb.WriteString(`    "github.com/wargasipil/protoc-gen-dispatcher/dispatch_core"` + "\n")
	sb.WriteString(`    ` + pkgImport + "\n")

	sb.WriteString(")\n\n")

	return sb.String()
}

func (g *Generator) GenerateServices() string {
	file := g.file

	var sb strings.Builder
	var funcb strings.Builder
	var stfield strings.Builder

	for _, service := range g.file.GetService() {
		name := service.GetName()
		structName := name + "Dispatcher"
		servicename := "/" + file.GetPackage() + "." + name + "/"

		// writing struct

		sb.WriteString("type " + name + "Dispatcher struct {\n")

		for _, method := range service.GetMethod() {
			if method.GetServerStreaming() {
				continue
			}

			if method.GetClientStreaming() {
				continue
			}

			procedureName := servicename + method.GetName()

			// writing struct
			inputType := g.getInputType(method.GetInputType())
			lowerName := strings.ToLower(method.GetName())
			sb.WriteString("    " + lowerName + " *dispatch_core.Client[" + inputType + "]\n")

			// writing func
			// func (s *HelloServiceDispatcher) Hello(ctx context.Context, req *connect.Request[*example.HelloRequest]) error {
			// 	return s.hello.CallUnary(ctx, req)
			// }
			funcb.WriteString("func (s *" + structName + ") " + method.GetName() + "(ctx context.Context, req *connect.Request[" + inputType + "]) error {\n")
			funcb.WriteString("\treturn s." + lowerName + ".CallUnary(ctx, req)\n")
			funcb.WriteString("}\n\n")

			// writing struct field
			// hello: dispatch_core.NewClientDispather[example.HelloRequest](exampleconnect.HelloServiceHelloProcedure, dispatcher, options...)
			stfield.WriteString("\t\t" + lowerName + ": dispatch_core.NewClientDispather[" + inputType + "](\"" + procedureName + "\", dispatcher, options...),\n")

		}
		sb.WriteString("}\n\n")

		// writing initiate
		sb.WriteString("func New" + structName + "(dispatcher dispatch_core.Dispatcher, options ...dispatch_core.ClientOption) *" + structName + " {\n")
		sb.WriteString("\treturn &" + structName + "{\n")
		sb.WriteString(stfield.String())
		sb.WriteString("\t}\n")
		sb.WriteString("}\n\n")
	}

	return sb.String() + funcb.String()
}

func (g *Generator) getInputType(d string) string {
	inputTypes := strings.Split(d, ".")
	return strings.Join(inputTypes[2:], ".")
}

func sanitizePackage(protoPkg string) string {
	if protoPkg == "" {
		return "proto"
	}
	parts := strings.Split(protoPkg, ".")
	return parts[len(parts)-2] + "_dispatcher"
}
