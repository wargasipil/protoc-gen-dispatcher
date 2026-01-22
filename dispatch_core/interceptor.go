package dispatch_core

import (
	"context"

	"connectrpc.com/connect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type sourceInterceptor struct{}

// WrapStreamingClient implements connect.Interceptor.
func (s *sourceInterceptor) WrapStreamingClient(connect.StreamingClientFunc) connect.StreamingClientFunc {
	panic("unimplemented")
}

// WrapStreamingHandler implements connect.Interceptor.
func (s *sourceInterceptor) WrapStreamingHandler(connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	panic("unimplemented")
}

// WrapUnary implements connect.Interceptor.
func (s *sourceInterceptor) WrapUnary(handler connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		// source, _ := custom_connect.GetRequestSource(ctx)
		// if source == nil {
		// 	return handler(ctx, req)
		// }

		// sourceString, err := custom_connect.RequestSourceSerialize(source)
		// if err != nil {
		// 	return nil, err
		// }

		// // req.Peer().Query.Set("x-pdc-source", sourceString)
		// req.Header().Set("X-Pdc-Source", sourceString)

		// if req.Header().Get("Authorization") == "" {
		// 	token, _ := custom_connect.GetAuthToken(ctx)
		// 	if token != "" {
		// 		req.Header().Set("Authorization", token)
		// 	}
		// }

		return handler(ctx, req)
	}
}

type telemetryInterceptor struct{}

// WrapStreamingClient implements connect.Interceptor.
func (t *telemetryInterceptor) WrapStreamingClient(connect.StreamingClientFunc) connect.StreamingClientFunc {
	panic("unimplemented")
}

// WrapStreamingHandler implements connect.Interceptor.
func (t *telemetryInterceptor) WrapStreamingHandler(connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	panic("unimplemented")
}

// WrapUnary implements connect.Interceptor.
func (t *telemetryInterceptor) WrapUnary(handler connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, ar connect.AnyRequest) (connect.AnyResponse, error) {
		span := trace.SpanFromContext(ctx)
		res, err := handler(ctx, ar)
		if err != nil {
			span.SetStatus(codes.Error, err.Error())
		}

		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(ar.Header()))

		return res, err
	}
}
