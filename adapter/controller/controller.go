package controller

import (
	"context"

	"github.com/bayu-aditya/ideagate/backend/client/worker-rest/config"
	pbController "github.com/bayu-aditya/ideagate/backend/model/gen-go/client/controller/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IControllerAdapter interface {
	GetListEndpoint(ctx context.Context) (*pbController.GetListEndpointResponse, error)
	GetWorkflow(ctx context.Context, entrypointID string) (*pbController.GetWorkflowResponse, error)
}

func New() (IControllerAdapter, error) {
	cfg := config.Get()

	conn, err := grpc.NewClient(cfg.Controller.Url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pbController.NewControllerServiceClient(conn)

	return &controllerAdapter{
		conn:   conn,
		client: client,
	}, nil
}

type controllerAdapter struct {
	conn   *grpc.ClientConn
	client pbController.ControllerServiceClient
}

func (c *controllerAdapter) GetListEndpoint(ctx context.Context) (*pbController.GetListEndpointResponse, error) {
	resp, err := c.client.GetListEndpoint(ctx, &pbController.GetListEndpointRequest{})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *controllerAdapter) GetWorkflow(ctx context.Context, entrypointID string) (*pbController.GetWorkflowResponse, error) {
	resp, err := c.client.GetWorkflow(ctx, &pbController.GetWorkflowRequest{
		EntrypointId: entrypointID,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *controllerAdapter) Close() error {
	return c.conn.Close()
}
