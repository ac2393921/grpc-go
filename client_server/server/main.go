package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/ac2393921/grpc-go/client_server/pb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const port = ":50051"

type ServerClientSide struct {
	pb.UnimplementedUploadServer
}

func (s *ServerClientSide) Upload(stream pb.Upload_UploadServer) error {
	var sum int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			message := fmt.Sprintf("Done: sum = %d", sum)
			return stream.SendAndClose(&pb.UploadReply{
				Message: message,
			})
		}
		if err != nil {
			return err
		}

		fmt.Println(req.GetValue())
		sum += req.GetValue()
	}
}

func set() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return errors.Wrap(err, "failed port")
	}
	s := grpc.NewServer()
	var server ServerClientSide
	pb.RegisterUploadServer(s, &server)
	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "failed serve")
	}

	return nil
}

func main() {
	fmt.Println("Start!")
	if err := set(); err != nil {
		log.Fatalf("%v", err)
	}
}
