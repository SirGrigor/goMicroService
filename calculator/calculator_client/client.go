package main

import (
	"context"
	"fmt"
	"goMicroService/calculator/calculatorpb"
	"google.golang.org/grpc"
	"log"
)

func main()  {
	fmt.Println("Caclulator CLient")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect := %v", err)
	}
	defer conn.Close()
	c := calculatorpb.NewCalculatorServiceClient(conn)
	doUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient)  {
	log.Println("Starting to do SUM unary")
	request := &calculatorpb.SumRequest{
		FirstNumber:  5,
		SecondNumber: 40,
	}
	res, err := c.Sum(context.Background(), request)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from greet: %v", res.SumResult)
}