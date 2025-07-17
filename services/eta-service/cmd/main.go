package cmd

import (
	"log"
	"net"

	pb "github.com/AlexanderZah/order-tracking/services/eta-service/gen/go/etaservice/v1"
	etaServer "github.com/AlexanderZah/order-tracking/services/eta-service/internal/handler/server"
	"google.golang.org/grpc"
)

func main() {
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	etaServ := etaServer.New()
	pb.RegisterETAServiceServer(grpcServer, etaServ)

	log.Println("gRPC server started on :50051")
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
