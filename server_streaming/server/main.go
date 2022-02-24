package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ac2393921/grpc-go/server_streaming/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const port = ":50051"

type ServerSide struct {
	pb.UnimplementedNotificationServer
}

func (s *ServerSide) Notification(req *pb.NotificationRequest, stream pb.Notification_NotificationServer) error {
	fmt.Println("recive request")
	for i := int32(0); i < req.GetNum(); i++ {
		message := fmt.Sprintf("%d", i)
		if err := stream.Send(&pb.NotificationReply{
			Message: message,
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}

func set() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return errors.Wrap(err, "failed port")
	}

	s := grpc.NewServer()
	var server ServerSide
	pb.RegisterNotificationServer(s, &server)
	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "failed start server")
	}
	return nil
}

func main() {
	fmt.Println("Start!")
	if err := set(); err != nil {
		log.Fatalf("%v", err)
	}
}
