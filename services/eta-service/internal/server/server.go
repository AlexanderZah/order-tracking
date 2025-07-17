package server

import (
	"context"
	"math/rand"

	pb "github.com/AlexanderZah/order-tracking/services/eta-service/gen/go/etaservice/v1"
)

type Server struct {
	pb.UnimplementedETAServiceServer
}

func New() *Server {
	return &Server{}
}

func (s *Server) GetETA(ctx context.Context, req *pb.Order) (*pb.ETAResponse, error) {
	eta := rand.Intn(21) + 10
	return &pb.ETAResponse{
		Eta: int32(eta),
	}, nil
}
