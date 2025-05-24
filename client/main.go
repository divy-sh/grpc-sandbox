package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/divy-sh/grpc-sandbox/sandbox"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("1: Unary\n2: Server Streaming\n3: Client Streaming\n4: Bidirectional streaming\n")
	input, _ := reader.ReadByte()
	switch input {
	case '1':
		UnaryClient()
	case '2':
		ServerStreamClient()
	case '3':
		ClientStreamClient()
	case '4':
		BidirectionalStreamClient()
	default:
		fmt.Println("invalid input")
	}
}

func UnaryClient() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewSandboxServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.UnaryCall(ctx, &pb.Request{Message: "Hello from client!"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Response from server: %s", resp.Reply)
}

func ServerStreamClient() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewSandboxServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := c.ServerStreamCall(ctx, &pb.Request{Message: "Hello server-stream!"})
	if err != nil {
		log.Fatalf("error calling ServerStreamCall: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			log.Println("Server stream ended")
			break
		}
		if err != nil {
			log.Fatalf("error receiving stream: %v", err)
		}
		log.Printf("Received from stream: %s", resp.Reply)
	}
}

func ClientStreamClient() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewSandboxServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := c.ClientStreamCall(ctx)
	if err != nil {
		log.Fatalf("could not start client stream: %v", err)
	}

	// Send multiple messages
	messages := []string{"Hello", "from", "the", "client", "stream!"}
	for _, msg := range messages {
		if err := stream.Send(&pb.Request{Message: msg}); err != nil {
			log.Fatalf("error sending message: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Close the stream and receive final response
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error receiving final response: %v", err)
	}
	log.Printf("Final response from server: %s", reply.Reply)
}

func BidirectionalStreamClient() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewSandboxServiceClient(conn)
	stream, err := client.BidiStreamCall(context.Background())
	if err != nil {
		log.Fatalf("error calling BidiStreamCall: %v", err)
	}

	// Send and receive in parallel
	done := make(chan struct{})

	// Receive messages from server
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error receiving: %v", err)
			}
			log.Printf("Received: %s", resp.Reply)
		}
		close(done)
	}()

	messages := []string{"Hey", "this", "is", "a", "bidirectional", "stream!"}
	for _, msg := range messages {
		log.Printf("Sending: %s", msg)
		if err := stream.Send(&pb.Request{Message: msg}); err != nil {
			log.Fatalf("error sending: %v", err)
		}
		time.Sleep(400 * time.Millisecond)
	}

	if err := stream.CloseSend(); err != nil {
		log.Fatalf("error closing send: %v", err)
	}

	<-done
}
