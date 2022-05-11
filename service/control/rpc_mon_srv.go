package control

import (
	"app"
	"context"
	"log"
	"net"

	pb "app/grpc/monitor"
	"google.golang.org/grpc"
)

// monitor.MonitorServer
type server struct {
	pb.UnimplementedMonitorServer
}

// Health implements monitor.MonitorServer
func (s *server) Health(ctx context.Context, in *pb.HealthData) (*pb.Response, error) {
	log.Printf("Received: %v", in.GetId())
	return &pb.Response{Code: 200, Message: "success"}, nil
}

// Notice implements monitor.MonitorServer
func (s *server) Notice(ctx context.Context, in *pb.BrokerData) (*pb.Response, error) {
	log.Printf("Received: %v", in.GetNumber())
	return &pb.Response{Code: 200, Message: "success"}, nil
}

func RunMonitorServer() {

	host := app.Yaml.Base.Monitor + ":" + app.Yaml.Base.PortRpc
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("monitor: failed to listen %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMonitorServer(s, &server{})
	log.Printf("monitor listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("monitor: failed to serve %v", err)
	}
}
