package control

import (
	"app"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"strconv"

	pb "app/grpc/monitor"
	"google.golang.org/grpc"
)

// monitor.MonitorServer
type server struct {
	pb.UnimplementedMonitorServer
}

// Health implements monitor.MonitorServer
func (s *server) Health(ctx context.Context, in *pb.HealthData) (*pb.Response, error) {
	log.Printf("Moniter: Received %v, %v, %v", in.GetId(), in.GetIp(), in.GetRole())
	return &pb.Response{Code: 200, Message: "success"}, nil
}

// Notice implements monitor.MonitorServer
func (s *server) Notice(ctx context.Context, in *pb.BrokerData) (*pb.Response, error) {
	log.Printf("Moniter: Received %v, %v", in.GetId(), in.GetNumber())
	return &pb.Response{Code: 200, Message: "success"}, nil
}

func RunMonitorServer() {

	p, _ := strconv.ParseInt(app.Yaml.Base.PortRpc, 10, 64)
	host := app.Yaml.Base.Monitor + ":" + strconv.FormatInt(p+1, 10)
	lis, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("monitor: failed to listen %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMonitorServer(s, &server{})
	fmt.Printf("monitor listening at %v\n", lis.Addr())
	log.Printf("monitor listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("monitor: failed to serve %v", err)
	}
}
