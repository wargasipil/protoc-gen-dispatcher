package main

import (
	"context"
	"log"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"connectrpc.com/connect"
	"github.com/wargasipil/protoc-gen-dispatcher/dispatch_core"
	"github.com/wargasipil/protoc-gen-dispatcher/example_gen/example/v1"
	"github.com/wargasipil/protoc-gen-dispatcher/example_gen/example/v1/example_dispatcher"
)

func main() {
	ctx := context.Background()

	client, err := cloudtasks.NewClient(ctx)
	if err != nil {
		log.Fatal(err)

	}

	defer client.Close()

	dispatcher := dispatch_core.NewGoogleDispatcher(client, "slow-event-task", "https://localhost:8080/")

	exam := example_dispatcher.NewHelloServiceDispatcher(dispatcher)
	err = exam.Hello(ctx, &connect.Request[example.HelloRequest]{
		Msg: &example.HelloRequest{
			Name: "asdasdasd",
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
