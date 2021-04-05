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
	"log"
	"net"

	pb "github.com/smiletrl/micro_ecommerce/service.product/external/rpc"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is rpc server for product
type server struct {
	pb.UnimplementedProductServer
}

func (s *server) GetProduct(ctx context.Context, in *pb.ProductID) (*pb.ProductDetail, error) {
	log.Printf("product id is: %s\n", in.Id)
	// get product from cache firstly and then db.
	// query db table directly
	return &pb.ProductDetail{Title: "mac", Amount: "1288", Stock: "12"}, nil
}

// build the server for product service.
// deploy to k8s.
func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterProductServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
