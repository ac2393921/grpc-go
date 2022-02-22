package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ac2393921/grpc-go/unary/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const port = ":50051"

type ServerUnary struct {
	pb.UnimplementedArithmeticServiceServer
}

func (s *ServerUnary) Addition(ctx context.Context, parameters *pb.Parameters) (*pb.Answer, error) {
	a := parameters.GetA()
	b := parameters.GetB()
	fmt.Println(a, b)
	answer := a + b

	return &pb.Answer{
		Value: answer,
	}, nil
}

func set() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return errors.Wrap(err, "failed port")
	}

	s := grpc.NewServer()
	var server ServerUnary
	pb.RegisterArithmeticServiceServer(s, &server)
	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "failed server start")
	}

	return nil
}

func main() {
	fmt.Println("Start")
	if err := set(); err != nil {
		log.Fatalf("%v", err)
	}
}
