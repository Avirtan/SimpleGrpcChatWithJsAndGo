package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/avirtan/ProtoExmpl/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var id int = 0

type client struct {
	stream pb.Greeter_StreamServer
	id     int
	done   chan error
}

type Server struct {
	pb.UnimplementedGreeterServer
	clients map[int]*client
	mu      sync.RWMutex
}

func (s *Server) SayHelloAgain(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.GetName()}, nil
}

func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello again " + in.GetName()}, nil
}

func (s *Server) ListResponse(req *pb.Request, stream pb.Greeter_ListResponseServer) error {
	responses := [...]*pb.Response{{Message: "1"}, {Message: "2"}, {Message: "3"}, {Message: "4"}}
	for _, res := range responses {
		if err := stream.Send(res); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) Stream(stream pb.Greeter_StreamServer) error {
	ctx := stream.Context()
	// headers, _ := metadata.FromIncomingContext(ctx)
	// tokenRaw := headers["authorization"]
	// if len(tokenRaw) == 0 {
	// 	fmt.Println("no token")
	// } else {
	// 	fmt.Println(tokenRaw[0])
	// }
	client := &client{
		id:     id,
		stream: stream,
		done:   make(chan error),
	}
	s.mu.Lock()
	s.clients[id] = client
	id++
	s.mu.Unlock()
	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				client.done <- err
			}
			fmt.Println(req)
			switch action := req.GetAction().(type) {
			case *pb.Request_TypeS:
				msg := action.TypeS.Mtypes
				resp := &pb.Response{Message: fmt.Sprintf("you: %v", action.TypeS.Mtypes)}
				if err := stream.Send(resp); err != nil {
					client.done <- err
				}
				s.broadCast(client.id, msg)
			case *pb.Request_TypeW:
				msg := action.TypeW.Mtypew
				resp := &pb.Response{Message: fmt.Sprintf("you %v", msg)}
				if err := stream.Send(resp); err != nil {
					client.done <- err
				}
				s.broadCast(client.id, msg)
			}
		}
	}()
	var doneError error
	select {
	case <-ctx.Done():
		doneError = ctx.Err()
	case doneError = <-client.done:
	}
	log.Printf(`stream done with error "%v"`, doneError)
	delete(s.clients, client.id)
	log.Printf("%v - removing client", client.id)

	return doneError
}

func (s *Server) broadCast(id int, msg string) {
	for _, client := range s.clients {
		if client.id != id {
			resp := &pb.Response{Message: fmt.Sprintf("%v: send %v", id, msg)}
			if err := client.stream.Send(resp); err != nil {
				client.done <- err
			}
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "192.168.1.185:30000")

	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	pb.RegisterGreeterServer(grpcServer, &Server{clients: make(map[int]*client)})
	fmt.Println("server run")
	grpcServer.Serve(listener)
}
