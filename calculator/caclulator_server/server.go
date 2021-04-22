package main

import (
	"context"
	"fmt"
	"goMicroService/calculator/calculatorpb"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type server struct{}

func (s *server) PrimeNumberDecomposition(request *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	log.Printf("Received PrimeNumberDecomposition: %v", request)
	inputNumber := request.GetPrimeNumber().PrimeNumber
	var k int32 = 2
	for inputNumber > 1 {
		if inputNumber%k == 0 && k < 10 {
			log.Println(k)
			inputNumber = inputNumber / k
			res := &calculatorpb.PrimeNumberDecompositionResponse{PrimeNumber: k}
			stream.Send(res)
		} else {
			k++
		}
	}
	return nil
}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Recieved SUM RPC: %v", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	sum := firstNumber + secondNumber
	res := &calculatorpb.SumResponse{SumResult: sum}
	return res, nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	log.Println("Calculate Average")
	sum := int64(0)
	totalIterations := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(totalIterations)
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		sum += req.GetPrimenumber()
		totalIterations++
	}
}

func main() {
	fmt.Println("Calculator server ")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
