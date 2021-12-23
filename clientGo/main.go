package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	pb "github.com/avirtan/ProtoExmpl/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func simpleClient() {

}

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("192.168.1.185:30000", opts...)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}

	defer conn.Close()
	client := pb.NewGreeterClient(conn)
	//? 1
	//
	// request := &pb.HelloRequest{
	// 	Name: "test",
	// }
	// response, err := client.SayHello(context.Background(), request)

	// if err != nil {
	// 	grpclog.Fatalf("fail to dial: %v", err)
	// }

	// fmt.Println(response.Message)
	//?
	//? 2
	// initialize a pb.Rectangle
	// stream, err := client.ListResponse(context.Background(), &pb.Request{Message: "hello"})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// for {
	// 	feature, err := stream.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
	// 	}
	// 	log.Println(feature)
	// }
	//?
	//? 3
	// header := metadata.New(map[string]string{"authorization": resp.Token})
	// ctx := metadata.NewOutgoingContext(context.Background(), header)
	stream, err := client.Stream(context.Background())
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			fmt.Printf("%s\n", in.Message)
		}
	}()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		request := &pb.Request_TypeS{&pb.TypeS{Mtypes: text}}
		if err := stream.Send(&pb.Request{Action: request}); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
	}
	// <-waitc
	// stream.CloseSend()
	//?
}
