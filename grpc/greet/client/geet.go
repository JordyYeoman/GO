package main

import (
	"context"
	pb "github.com/JordyYeoman/GO/grpc/greet/proto"
	"log"
)

func doGreet(c pb.GreetServiceServer) {
	log.Println("doGreet was invoked")
	res, err := c.Greet(context.Background(), &pb.GreetRequest{
		FirstName: "Jordy",
	})

	if err != nil {
		log.Fatalf("Could not greet: %v\n", err)
	}

	log.Printf("Greeting: %s\n", res.Result)
}
