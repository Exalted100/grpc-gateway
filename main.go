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

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"

	// "fmt"
	"log"
	// "net"
	"net/http"

	pb "github.com/Exalted100/go-grpc/helloworld"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:8000", "gRPC server endpoint")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	flag.Parse()
	// lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	// if err != nil {
	// 	log.Fatalf("failed to listen: %v", err)
	// }

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	// s := grpc.NewServer()
	// pb.RegisterGreeterServer(s, &server{})
	pb.RegisterGreeterHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	// fmt.Println(*grpcServerEndpoint)

	err := mux.HandlePath("GET", "/hello/{name}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		w.Write([]byte("hello " + pathParams["name"]))
	})
	if err != nil {
		panic(err)
	}

	// log.Printf("server listening at %v", lis.Addr())
	// if err := s.Serve(lis); err != nil {
	// 	log.Fatalf("failed to serve: %v", err)
	// }
	http.ListenAndServe(":8081", mux)
}

// func main() {
// 	// creating mux for gRPC gateway. This will multiplex or route request different gRPC service
// 	mux := runtime.NewServeMux()
// 	// setting up a dail up for gRPC service by specifying endpoint/target url
// 	err := pb.RegisterGreeterHandlerFromEndpoint(context.Background(), mux, "localhost:8080", []grpc.DialOption{grpc.WithInsecure()})

// 	err = mux.HandlePath("GET", "/hello/{name}", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
// 		w.Write([]byte("hello " + pathParams["name"]))
// 	})
// 	if err != nil {
// 		panic(err)
// 	}

// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// Creating a normal HTTP server
// 	server := http.Server{
// 		Handler: mux,
// 	}
// 	// creating a listener for server
// 	l, err := net.Listen("tcp", ":8081")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	// start server
// 	err = server.Serve(l)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
