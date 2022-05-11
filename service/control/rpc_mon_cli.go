package control

import (
	"app"
	"context"
	"log"
	"time"

	pb "app/grpc/monitor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func MonitorClient() {
	// Set up a connection to the server.
	host := app.Yaml.Base.Monitor + ":" + app.Yaml.Base.PortRpc
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
