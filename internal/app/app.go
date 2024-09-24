package app

import (
	"github.com/BazaarTrade/ApiGatewayService/internal/api/rest"
	ws "github.com/BazaarTrade/ApiGatewayService/internal/api/websocket"
	"github.com/BazaarTrade/GeneratedProto/pb"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

func Run() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return
	}

	pbClient := pb.NewMatchingEngineClient(conn)
	e := echo.New()

	rest.NewHandler(ws.NewHub([]string{"OrderUpdate"}), pbClient).Init(e)
	e.Start(":8080")
}
