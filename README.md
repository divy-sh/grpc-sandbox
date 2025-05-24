
# grpc-sandbox

**grpc-sandbox** is a project designed for testing and experimenting with gRPC services. It provides a clean structure for setting up example services and clients, allowing for rapid prototyping, debugging, and performance validation.

## Features

- Example gRPC server and client implementations  
- Covers all RPC types: unary, server streaming, client streaming, and bidirectional streaming  
- Test scenarios for integration and behavior validation  
- Simple logging and debugging hooks  
- Optional Docker setup for containerized testing  

## Project Structure

grpc-test-lab/  
├── proto/             # Protocol Buffer definitions  
├── server/            # gRPC server code  
├── client/            # gRPC client code  
├── tests/             # Test scripts and assertions  
├── docker/            # Docker and Docker Compose setup  
├── Makefile           # Common CLI commands  
└── README.md  

## Requirements

- A supported language environment (e.g. Go, Python, Node.js, etc.)  
- protoc Protocol Buffers compiler  
- gRPC library for the chosen language  
- Docker (optional, for containerized environments)  

## Setup

### Clone the repository

git clone https://github.com/divy-sh/grpc-sandbox.git
cd grpc-sandboc

### Compile Protocol Buffers (example for Go)

protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/*.proto

### Run the server

go run server/main.go

### Run the client

go run client/main.go

### Run tests

go test ./...


## License

MIT
