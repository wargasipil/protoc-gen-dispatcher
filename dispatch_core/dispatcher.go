package connect_dispatcher

import (
	"context"

	"connectrpc.com/connect"
)

type clientConfig struct {
	// URL         *url.URL
	endpoint    string
	Procedure   string
	Interceptor connect.Interceptor
}

type ClientOption interface {
	applyToClient(*clientConfig)
}

type chain struct {
	interceptors []connect.Interceptor
}

func (c *chain) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	for _, interceptor := range c.interceptors {
		next = interceptor.WrapUnary(next)
	}
	return next
}

func (c *chain) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next
}

func (c *chain) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next
}

type interceptorsOption struct {
	Interceptors []connect.Interceptor
}

func (o *interceptorsOption) chainWith(current connect.Interceptor) connect.Interceptor {
	if len(o.Interceptors) == 0 {
		return current
	}
	if current == nil && len(o.Interceptors) == 1 {
		return o.Interceptors[0]
	}
	if current == nil && len(o.Interceptors) > 1 {
		return newChain(o.Interceptors)
	}
	return newChain(append([]connect.Interceptor{current}, o.Interceptors...))
}

// applyToClient implements connect.Option.
func (i *interceptorsOption) applyToClient(config *clientConfig) {
	config.Interceptor = i.chainWith(config.Interceptor)
}

func WithInterceptors(interceptors ...connect.Interceptor) ClientOption {
	return &interceptorsOption{interceptors}
}

func newChain(interceptors []connect.Interceptor) *chain {
	// We usually wrap in reverse order to have the first interceptor from
	// the slice act first. Rather than doing this dance repeatedly, reverse the
	// interceptor order now.
	var chain chain
	for i := len(interceptors) - 1; i >= 0; i-- {
		if interceptor := interceptors[i]; interceptor != nil {
			chain.interceptors = append(chain.interceptors, interceptor)
		}
	}
	return &chain
}

func newClientConfig(endpoint string, options ...ClientOption) (*clientConfig, *connect.Error) {

	config := clientConfig{
		// URL:              endpoint,
		endpoint: endpoint,
		// Procedure:        protoPath,
	}

	for _, opt := range options {
		opt.applyToClient(&config)
	}

	return &config, nil
}

type Client[Req any] struct {
	config    *clientConfig
	Procedure string
	CallUnary func(ctx context.Context, r *connect.Request[Req]) error
	err       error
}

type Dispatcher func(ctx context.Context, procedure string, r connect.AnyRequest) error

func NewClientDispather[Req any](procedure string, dispatcher Dispatcher, options ...ClientOption) *Client[Req] {
	client := &Client[Req]{}

	options = append(
		options,
		WithInterceptors(&sourceInterceptor{}, &telemetryInterceptor{}),
	)

	config, err := newClientConfig(
		procedure,
		options...)

	if err != nil {
		client.err = err
		return client
	}
	client.config = config
	client.Procedure = procedure

	var unaryFunc connect.UnaryFunc = func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {

		return nil, dispatcher(ctx, procedure, ar)
	}

	if interceptor := config.Interceptor; interceptor != nil {
		// interceptor is the full chain of all interceptors provided
		unaryFunc = interceptor.WrapUnary(unaryFunc)
	}

	client.CallUnary = func(ctx context.Context, r *connect.Request[Req]) error {
		_, err := unaryFunc(ctx, r)
		return err
	}

	return client
}
