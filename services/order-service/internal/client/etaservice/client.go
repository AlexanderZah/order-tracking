package etaservice

import (
	"context"

	pb "github.com/AlexanderZah/order-tracking/services/eta-service/gen/go/etaservice/v1"
	"google.golang.org/grpc"
)

type Client struct {
	eta pb.ETAServiceClient
}

func New(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	cli := pb.NewETAServiceClient(conn)
	return &Client{eta: cli}, nil
}

func (c *Client) GetETA(ctx context.Context, addr string) (int32, error) {
	resp, err := c.eta.GetETA(ctx, &pb.Order{
		DeliveryAddress: addr,
	})

	if err != nil {
		return 0, err
	}

	return resp.GetEta(), nil
}
