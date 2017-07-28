package server

import (
	"fmt"
	kexec "kola/exec"
	pb "kola/pb"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const (
	port = 5051
)

//KolaServer This is the server accept the client dial
type KolaServer struct{}

//Get the call information send from client
func (k *KolaServer) Get(ctx context.Context, in *pb.KolaRequest) (*pb.KolaReply, error) {
	message := in.Key
	c := kexec.NewCmd("echo", "'test'")
	statusChan := c.Start()
	finalStatus := <-statusChan
	return &pb.KolaReply{Props: append(message, finalStatus.Stdout...)}, nil
}

//StartServer used to start the Server
func StartServer() {
	lis, err := net.Listen("tcp4", fmt.Sprintf(":%d", port))
	if err != nil {
		grpclog.Fatalf("failed to listen: %v \n", err)
	}
	kolaServer := grpc.NewServer()
	pb.RegisterKolaAgentServer(kolaServer, &KolaServer{})
	if err := kolaServer.Serve(lis); err != nil {
		grpclog.Fatalf("failed to server %v \n", err)
	}
}
