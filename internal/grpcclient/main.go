package main

import (
	"log"

	pb "github.com/sik0-o/gorcon-restarter/v2/internal/proto/grpcrcon"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("not to able connect:%v", err)
	}

	defer conn.Close()

	c := pb.NewGRPCRCONServiceClient(conn)
	announceServer(c)
}
