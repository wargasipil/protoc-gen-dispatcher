package example_dispatcher

import (
    "context"
    connect "connectrpc.com/connect"
    "github.com/wargasipil/protoc-gen-dispatcher/dispatch_core"
    v1 "github.com/wargasipil/protoc-gen-dispatcher/example_gen/example/v1"
)

type HelloServiceDispatcher struct {
    hello *dispatch_core.Client[v1.HelloRequest]
}

func NewHelloServiceDispatcher(dispatcher dispatch_core.Dispatcher, options ...dispatch_core.ClientOption) *HelloServiceDispatcher {
	return &HelloServiceDispatcher{
		hello: dispatch_core.NewClientDispather[v1.HelloRequest]("/example.v1.HelloService/Hello", dispatcher, options...),
	}
}

func (s *HelloServiceDispatcher) Hello(ctx context.Context, req *connect.Request[v1.HelloRequest]) error {
	return s.hello.CallUnary(ctx, req)
}

