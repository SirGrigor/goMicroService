package main

import (
	"context"
	"fmt"
	"goMicroService/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func doBiDiStreaming(c calculatorpb.CalculatorServiceClient) {
	log.Println("Starting to do a FindMax BiDi Streaming Rpc")
	stream, err := c.FindMaximum(context.Background())

	if err != nil {
		log.Fatalf("Error while opening stream and calling FindMaximum: %v", err)
	}
	waitc := make(chan struct{})

	//send
	go func() {
		numbers := []int64{4, 4, 5, 6, 7, 8, 91}
		for _, number := range numbers {
			fmt.Printf("Number sent: %v", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				PrimeNumber: number,
			})
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()

	}()
	//receive

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Problem while reading: %v", err)
			}
			maximum := res.GetMaxNumber()
			fmt.Printf("Received a new maximum of...%v", maximum)
		}
		close(waitc)
	}()
	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Start to do a SquareRoot Unary RPC")

	//correct call
	doErrorCall(c, 10)

	//error call
	doErrorCall(c, -2)

}

func doErrorCall(c calculatorpb.CalculatorServiceClient, number int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: number})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			//user ok
			fmt.Printf("Error message from server: %v", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We sent a negative number!")
				return
			}
		} else {
			log.Fatalf("Big error calling Square rootL %v", err)
			return
		}
	}
	fmt.Printf("Result of square root of %v: %v", number, res.GetNumberRoot())
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
	//doAverageCalculation(c)
	//doBiDiStreaming(c)
	doErrorUnary(c)
}
