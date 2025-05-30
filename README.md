
# grpc-sandbox

**grpc-sandbox**

## Features

- Example gRPC server and client implementations  
- Covers all RPC types: unary, server streaming, client streaming, and bidirectional streaming  
- Test scenarios for integration and behavior validation  
- Simple logging and debugging hooks  
- Optional Docker setup for containerized testing  

## Project Structure

grpc-sandbox/  
├── proto/             # Protocol Buffer definitions  
├── server/            # gRPC server code  
├── client/            # gRPC client code  
├── tests/             # Test scripts and assertions  
└── README.md

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
