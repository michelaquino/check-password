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
	// echoServer.GET("/", Index)
	echoServer.GET("/check", handler.GetCheck)
	echoServer.POST("/check", handler.PostCheck)

	echoServer.Logger.Fatal(echoServer.Start(":8888"))
}
