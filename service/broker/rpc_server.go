package broker

import (
	"app"
	"context"
	"log"
	"net"

	pb "app/grpc/forward"
	"google.golang.org/grpc"
)

// server is used to implement forward.ForwardingServer.
type server struct {
	pb.UnimplementedForwardingServer
}

// SendMessage implements forward.ForwardingServer
func (s *server) SendMessage(ctx context.Context, in *pb.MsgRequest) (*pb.MsgReply, error) {
	log.Printf("Received: %v", in.GetPayload())
	return &pb.MsgReply{Code: 200, Message: "success"}, nil
}

func StartRPC() {

	// FIXME:: host and port
	lis, err := net.Listen("tcp", app.Yaml.Base.RpcHost)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterForwardingServer(s, &server{})
	log.Printf("rpc listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
