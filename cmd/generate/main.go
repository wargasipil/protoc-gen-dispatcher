package main

import (
	"io"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(input, req); err != nil {
		panic(err)
	}

	// 2. Prepare CodeGeneratorResponse
	res := &pluginpb.CodeGeneratorResponse{}

	for _, file := range req.GetProtoFile() {
		if !contains(req.GetFileToGenerate(), file.GetName()) {
			continue
		}

		gen := Generator{
			file: file,
		}

		if !gen.IsHaveService() {
			continue
		}

		// mockFile := GenerateMockFile(file)
		mockFile := gen.Generate()
		res.File = append(res.File, mockFile)
	}

	// 3. Write CodeGeneratorResponse to stdout
	out, err := proto.Marshal(res)
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(out)
}

func contains(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}
