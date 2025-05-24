package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/divy-sh/grpc-sandbox/sandbox"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedSandboxServiceServer
}

func main() {
	go RunUnaryServer()
	go RunServerStreamServer()
	go RunClientStreamServer()
	RunBidirectionalStreamServer()
}

func (s *server) ClientStreamCall(stream pb.SandboxService_ClientStreamCallServer) error {
	var messages []string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			reply := fmt.Sprintf("Received %d messages: %v", len(messages), messages)
			return stream.SendAndClose(&pb.Response{Reply: reply})
		}
		if err != nil {
			return err
		}
		log.Printf("Received from client: %s", req.Message)
		messages = append(messages, req.Message)
	}
}

func RunClientStreamServer() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSandboxServiceServer(s, &server{})
	log.Println("Server listening on :50053")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) ServerStreamCall(req *pb.Request, stream pb.SandboxService_ServerStreamCallServer) error {
	log.Printf("Received message: %s", req.Message)

	// Simulate sending 5 messages back to the client
	for i := 0; i < 5; i++ {
		resp := &pb.Response{Reply: fmt.Sprintf("Response %d to message: %s", i+1, req.Message)}
		if err := stream.Send(resp); err != nil {
			return err
		}
		time.Sleep(time.Second) // simulate delay
	}

	return nil
}

func RunServerStreamServer() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSandboxServiceServer(s, &server{})
	log.Println("Server listening on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) UnaryCall(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Printf("Received message: %s", req.Message)
	return &pb.Response{Reply: "Hello from server, you said: " + req.Message}, nil
}

func RunUnaryServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSandboxServiceServer(s, &server{})
	log.Println("Server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) BidiStreamCall(stream pb.SandboxService_BidiStreamCallServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("Client stream ended")
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("Received from client: %s", req.Message)

		// Respond immediately to each message
		resp := &pb.Response{
			Reply: fmt.Sprintf("Echo: %s", req.Message),
		}
		if err := stream.Send(resp); err != nil {
			return err
		}
		time.Sleep(500 * time.Millisecond) // Simulate processing time
	}
}

func RunBidirectionalStreamServer() {
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSandboxServiceServer(s, &server{})
	log.Println("Server listening on :50054")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
