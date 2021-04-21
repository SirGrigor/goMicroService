package main

import (
	"context"
	"fmt"
	"goMicroService/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with, %v", req)
	firstname := req.GetGreeting().GetFirstName()

	results := "hello " + firstname
	res := &greetpb.GreetResponse{Result: results}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	log.Println("Greet Many time function Invoked")
	FirstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + FirstName + " number " + strconv.Itoa(i)
		res := &greetpb.GreatManyTimesResponse{Result: result}
		stream.Send(res)
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	fmt.Println("Hello ")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
