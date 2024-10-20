package main

import (
	"time"

	"github.com/BazaarTrade/ApiGatewayService/internal/app"
)

func main() {
	time.Sleep(time.Second * 5)
	app.Run()
}
