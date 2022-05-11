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

	host := app.Yaml.Base.Host + ":" + app.Yaml.Base.PortRpc
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("broker: failed to listen %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterForwardingServer(s, &server{})
	log.Printf("broker listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("broker: failed to serve %v", err)
	}
}
