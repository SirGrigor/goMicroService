package main

import (
	"context"
	"fmt"
	"goMicroService/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func doPrimeNumberDecompositionStreaming(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting to Stream NumberDecomposition RPC...")
	reg := &calculatorpb.PrimeNumberDecompositionRequest{
		PrimeNumber: &calculatorpb.PrimeNumber{
			PrimeNumber: 120,
		},
	}
	resStream, err := c.PrimeNumberDecomposition(context.Background(), reg)
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
		log.Printf("Response from PrimeNumberDecomposition %v", msg.GetPrimeNumber())
	}
}
func doUnary(c calculatorpb.CalculatorServiceClient) {
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

func doAverageCalculation(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting to do Average Calculation")
	numbers := []int32{3, 5, 6, 7, 8}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling Compute Average %v", err)
	}

	for _, n := range numbers {
		log.Printf("Send number: %v", n)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Primenumber: int64(n),
		})
		time.Sleep(1 * time.Second)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while recive from Average: %v", err)
	}
	log.Printf("Average: %v", res)
}

func main() {
	fmt.Println("Caclulator CLient")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect := %v", err)
	}
	defer conn.Close()
	c := calculatorpb.NewCalculatorServiceClient(conn)
	//doUnary(c)
	//doPrimeNumberDecompositionStreaming(c)
	doAverageCalculation(c)
}
