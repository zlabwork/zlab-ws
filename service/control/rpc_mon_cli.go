package control

import (
	"app"
	"context"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"

	pb "app/grpc/monitor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func monitorClientDemo() {
	// Set up a connection to the server.
	p, _ := strconv.ParseInt(app.Yaml.Base.PortRpc, 10, 64)
	host := app.Yaml.Base.Monitor + ":" + strconv.FormatInt(p+1, 10)
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewMonitorClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// TODO: some client data
	r, err := c.Health(ctx, &pb.HealthData{Role: "broker", Id: 1})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Result: %d, %s", r.GetCode(), r.GetMessage())
}

func MonitorConn() (*grpc.ClientConn, error) {

	// Set up a connection to the server.
	p, _ := strconv.ParseInt(app.Yaml.Base.PortRpc, 10, 64)
	host := app.Yaml.Base.Monitor + ":" + strconv.FormatInt(p+1, 10)
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return nil, err
	}
	return conn, nil
}
