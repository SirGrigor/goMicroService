package main

import (
	"context"
	"fmt"
	"goMicroService/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	reg := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ilja",
			LastName:  "Grigorjev",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), reg)
	if err != nil {
		log.Fatalf("Error while calling %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes %v", msg.GetResult())
	}
}

func doUnary(c greetpb.GreetServiceClient) {
	log.Println("Starting to do unary")
	request := &greetpb.GreetRequest{Greeting: &greetpb.Greeting{
		FirstName: "Ilja",
		LastName:  "Grigorjev",
	},
	}
	res, err := c.Greet(context.Background(), request)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from greet: %v", res.Result)
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	log.Println("Starting to do a Client Streaming RPC...")
	requests := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Ilja",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Ilja2",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Ilja3",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Ilja4",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Ilja5",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Ilja6",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreetL %v", err)
	}

	for _, req := range requests {
		log.Printf("Senging Request: %v", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while reciveing response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet Response: %v", res)
}

func main() {
	fmt.Println("Hello, Iam a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect := %v", err)
	}
	defer conn.Close()
	c := greetpb.NewGreetServiceClient(conn)
	//doUnary(c)
	//doServerStreaming(c)
	doClientStreaming(c)
}
