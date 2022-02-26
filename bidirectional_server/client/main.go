package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/ac2393921/grpc-go/bidirectional_server/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func receive(stream pb.Chat_ChatClient) error {
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			log.Printf("From Server: %s", in.Message)

			stream.Send(&pb.ChatRequest{
				Message: time.Now().Format("2022-02-26 00:00:00"),
			})
		}
	}()
	<-waitc
	return nil
}

func request(stream pb.Chat_ChatClient) error {
	return stream.Send(&pb.ChatRequest{
		Message: "こんにちは",
	})
}

func chat(client pb.ChatClient) error {
	stream, err := client.Chat(context.Background())
	if err != nil {
		return err
	}
	if err := request(stream); err != nil {
		return err
	}
	if err := receive(stream); err != nil {
		return err
	}
	stream.CloseSend()
	return nil
}

func exec() error {
	address := "localhost:50051"
	conn, err := grpc.Dial(
		address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		return errors.Wrap(err, "connection error")
	}
	defer conn.Close()
	client := pb.NewChatClient(conn)
	return chat(client)
}

func main() {
	if err := exec(); err != nil {
		log.Println(err)
	}
}
