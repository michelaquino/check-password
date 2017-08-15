package main

import (
	"gitlab.globoi.com/michel.aquino/check-password/handler"
	"gitlab.globoi.com/michel.aquino/check-password/templates"

	"github.com/labstack/echo"
)

func main() {
	echoServer := echo.New()

	// Static files
	echoServer.Static("/static", "static")

	echoServer.Renderer = templates.ViewTemplates
	echoServer.GET("/", handler.GetListCredentials)
	echoServer.GET("/credentials", handler.GetCredentials)
	echoServer.POST("/credentials", handler.PostCredentials)

	echoServer.Logger.Fatal(echoServer.Start(":8888"))
}
