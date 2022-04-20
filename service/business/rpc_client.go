package business

import (
	"app"
	"context"
	"log"
	"time"

	pb "app/grpc/forward"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartRPC() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(app.Yaml.Base.RpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewForwardingClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SendMessage(ctx, &pb.MsgRequest{Payload: []byte("test 123456")})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Result: %d, %s", r.GetCode(), r.GetMessage())
}
