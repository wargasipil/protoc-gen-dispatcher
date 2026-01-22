package connect_dispatcher

import (
	"context"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"connectrpc.com/connect"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func NewGoogleDispatcher(
	client *cloudtasks.Client,
	queue string,
	endpoint string,
) Dispatcher {
	return func(ctx context.Context, procedure string, req connect.AnyRequest) error {

		content, err := protojson.Marshal(req.Any().(proto.Message))
		if err != nil {
			return err
		}

		headers := req.Header()
		hh := propagation.HeaderCarrier(headers)
		otel.GetTextMapPropagator().Inject(ctx, hh)

		reqheaders := make(map[string]string)
		for k, v := range headers {
			if len(v) > 0 {
				reqheaders[k] = v[0]
			}
		}

		// reqheaders["Content-Type"] = "application/grpc-web"
		reqheaders["Content-Type"] = "application/json"
		reqheaders["Connect-Protocol-Version"] = "1"

		httpreq := &cloudtaskspb.Task_HttpRequest{
			HttpRequest: &cloudtaskspb.HttpRequest{
				Url:        endpoint + procedure,
				HttpMethod: cloudtaskspb.HttpMethod_POST,
				Headers:    reqheaders,
				Body:       content,
			},
		}

		task := cloudtaskspb.CreateTaskRequest{
			Parent: queue,
			Task: &cloudtaskspb.Task{
				MessageType: httpreq,
			},
		}
		_, err = client.CreateTask(ctx, &task)

		return err
	}
}
