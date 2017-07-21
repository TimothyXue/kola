package client

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	pb "kola/pb"
)

const (
	serverAddr = "localhost:5051"
)

func StartClient() {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalf("did not connect： %v \n", err)
	}
	defer conn.Close()
	c := pb.NewKolaAgentClient(conn)

	//now try to talk with the server
	if r, err := c.Get(context.Background(), &pb.KolaRequest{Key: "test"}); err != nil {
		grpclog.Fatalf("could not get message from kola server %v \n", err)
	} else {
		grpclog.Infof("Returned Messages as %v \n", r.Props)
	}
}
