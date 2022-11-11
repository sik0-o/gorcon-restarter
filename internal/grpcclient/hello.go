package main

import (
	"context"
	"log"

	pb "github.com/sik0-o/gorcon-restarter/v2/internal/proto/grpcrcon"
)

func announceServer(c pb.GRPCRCONServiceClient) {
	res, err := c.Servers(context.Background(), &pb.ServersRequest{
		Name: "sikO_o",
	})

	if err != nil {
		log.Fatalf("unable to get response %v", err)
	}
	log.Printf("Greeting: %v", res.Greeting)

	if len(res.Servers) <= 0 {
		return
	}

	c.Announce(context.Background(), &pb.AnnounceRequest{
		ServerName: res.Servers[0].Name,
		Announce:   "My fucking announce for you bitches",
	})
}

func doGreet(c pb.GRPCRCONServiceClient) {
	res, err := c.Servers(context.Background(), &pb.ServersRequest{
		Name: "sikO_o",
	})

	if err != nil {
		log.Fatalf("unable to get response %v", err)
	}
	log.Printf("Greeting: %v", res.Greeting)
	log.Printf("Servers: %v", res.Servers)
}
