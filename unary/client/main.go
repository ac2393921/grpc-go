package main

import (
	"context"
	"log"
	"time"

	"github.com/ac2393921/grpc-go/unary/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func request(client pb.ArithmeticServiceClient, a, b int32) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()
	Parameters := pb.Parameters{
		A: a,
		B: b,
	}
	answer, err := client.Addition(ctx, &Parameters)
	if err != nil {
		return errors.Wrap(err, "failed reply")
	}
	log.Printf("Answer: %d", answer.GetValue())
	return nil
}

func addition(a, b int32) error {
	address := "localhost:50051"
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return errors.Wrap(err, "failed connection")
	}
	defer conn.Close()
	client := pb.NewArithmeticServiceClient(conn)
	return request(client, a, b)
}

func main() {
	a := int32(300)
	b := int32(500)
	if err := addition(a, b); err != nil {
		log.Fatalf("%v", err)
	}
}
