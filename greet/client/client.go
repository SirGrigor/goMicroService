package main

import (
	"context"
	"fmt"
	"goMicroService/greet/greetpb"
	"google.golang.org/grpc"
	"io"
	"log"
)

func main()  {
	fmt.Println("Hello, Iam a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect := %v", err)
	}
	defer conn.Close()
	c := greetpb.NewGreetServiceClient(conn)
	//doUnary(c)
	doServerStreaming(c)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")
	reg := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Ilja",
			LastName:  "Grigorjev",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), reg)
	if err != nil{
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

func doUnary(c greetpb.GreetServiceClient)  {
	log.Println("Starting to do unary")
	request := &greetpb.GreetRequest{ Greeting: &greetpb.Greeting{
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