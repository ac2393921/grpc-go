package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/ac2393921/grpc-go/bidirectional_server/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const port = ":50051"

type ServerBidirectional struct {
	pb.UnimplementedChatServer
}

func request(stream pb.Chat_ChatServer, message string) error {
	reply := fmt.Sprintf("Receive OK: ", message)
	return stream.Send(&pb.ChatReply{
		Message: reply,
	})
}

func (s *ServerBidirectional) Chat(stream pb.Chat_ChatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		message := in.GetMessage()
		fmt.Println("Receive: ", message)

		if err := request(stream, message); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
}

func set() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return errors.Wrap(err, "failed port")
	}
	s := grpc.NewServer()
	var server ServerBidirectional
	pb.RegisterChatServer(s, &server)
	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "failed server")
	}
	return nil
}

func main() {
	fmt.Println("Start")
	if err := set(); err != nil {
		log.Fatalf("%v", err)
	}
}
