package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/dim-pep/gRPC-CLient/interceptor"
	"github.com/dim-pep/gRPC-CLient/trading"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithUnaryInterceptor(interceptor.XRequestIDUnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("cannot connect to OrderService: %v", err)
	}
	defer conn.Close()

	client := trading.NewOrderServiceClient(conn)

	req := &trading.CreateOrderRequest{
		UserId:    "user-1",
		MarketId:  "BTC-USD",
		OrderType: "test",
		Price:     50000,
		Quantity:  1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.CreateOrder(ctx, req)
	if err != nil {
		log.Fatalf("CreateOrder failed: %v", err)
	}

	fmt.Println("Order ID:", resp.OrderId)
	fmt.Println("Статус:", resp.Status)

	getReq := &trading.GetOrderStatusRequest{
		OrderId: resp.OrderId,
		UserId:  "user-1",
	}

	statusResp, err := client.GetOrderStatus(ctx, getReq)
	if err != nil {
		log.Fatalf("GetOrderStatus failed: %v", err)
	}

	fmt.Println("Статус из GetOrderStatus:", statusResp.Status)
}
