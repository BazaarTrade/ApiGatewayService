package main

import (
	api "github.com/BazaarTrade/ApiGatewayService/internal/api/http"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	api.NewHandler().Init(e)
}
