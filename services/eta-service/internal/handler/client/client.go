package client

import (
	"context"
	"time"

	pb "github.com/AlexanderZah/order-tracking/services/eta-service/gen/go/etaservice/v1"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.ETAServiceClient
}

func New(address string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: pb.NewETAServiceClient(conn),
	}, nil
}
