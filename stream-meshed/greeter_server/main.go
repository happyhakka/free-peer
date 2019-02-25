/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

package main

import (
	"fmt"
	pb "github.com/chainforce/free-peer/stream-meshed/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"os"
)

var (
	port = ":50051"
	podName = os.Getenv("POD_NAME")
)

// server is used to implement helloworld.GreeterServer.
type server struct{
	Name string
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(stream pb.Greeter_SayHelloServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			log.Printf("EOF received on stream")
			continue
		}
		if err != nil {
			log.Fatal(fmt.Sprintf("error receiving on stream: %v", err))
		}
		log.Printf("Chaincode received: %v", in.Message)
		message := &pb.HelloMessage{Message: fmt.Sprintf("Repsonse to request: %v | from Chaincode: %v", in.Message, s.Name)}
		err = stream.Send(message)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	server := newServer()
	pb.RegisterGreeterServer(s, server)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func newServer() *server {
	s := &server{Name: podName}
	return s
}

