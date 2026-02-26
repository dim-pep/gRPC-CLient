package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dim-pep/gRPC-CLient/generated"
	"github.com/dim-pep/gRPC-CLient/interceptor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:8030",
		grpc.WithUnaryInterceptor(interceptor.XRequestIDUnaryClientInterceptor()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("cannot connect to OrderService: %v", err)
	}
	defer conn.Close()

	client := generated.NewOrderServiceClient(conn)

	req := &generated.CreateOrderRequest{
		UserId:    "user-1",
		MarketId:  "BTC-USD",
		OrderType: "test",
		Price:     50000,
		Quantity:  1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := client.CreateOrder(ctx, req)
	if err != nil {
		log.Fatalf("CreateOrder failed: %v", err)
	}

	fmt.Println("Order ID:", resp.OrderId)
	fmt.Println("Статус:", resp.Status)

	getReq := &generated.GetOrderStatusRequest{
		OrderId: resp.OrderId,
		UserId:  "user-1",
	}

	statusResp, err := client.GetOrderStatus(ctx, getReq)
	if err != nil {
		log.Fatalf("GetOrderStatus failed: %v", err)
	}

	fmt.Println("Статус из GetOrderStatus:", statusResp.Status)
}
