// Komunikacija po protokolu gRPC
// odjemalec

package main

import (
	"api/grpc/protobufStorage"
	"context"
	"fmt"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Client_subscriber(url string) {
	// vzpostavimo povezavo s strežnikom
	fmt.Printf("gRPC client connecting to %v\n", url)
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// vzpostavimo izvajalno okolje
	contextCRUD, cancel := context.WithCancel(context.Background())
	// contextCRUD, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// vzpostavimo vmesnik gRPC
	grpcClient := protobufStorage.NewCRUDClient(conn)

	// subscribamo in čakamo
	if stream, err := grpcClient.Subscribe(contextCRUD, &emptypb.Empty{}); err == nil {
		for {
			todoEvent, err := stream.Recv()
			if err == io.EOF {
				continue
			}
			fmt.Println(todoEvent.Action, "->", todoEvent.T)
		}
	} else {
		panic(err)
	}

}
