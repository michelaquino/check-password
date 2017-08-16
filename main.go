package main

import (
	"gitlab.globoi.com/michel.aquino/check-password/context"
	"gitlab.globoi.com/michel.aquino/check-password/handler"
	"gitlab.globoi.com/michel.aquino/check-password/templates"

	"github.com/labstack/echo"
)

func main() {
	context.GetAPIConfig()

	echoServer := echo.New()

	// Static files
	echoServer.Static("/static", "static")

	echoServer.Renderer = templates.ViewTemplates
	echoServer.GET("/", handler.GetListCredentials)
	echoServer.GET("/credentials", handler.GetCredentials)
	echoServer.POST("/credentials", handler.PostCredentials)

	echoServer.Logger.Fatal(echoServer.Start(":8888"))
}

// db.credentials.find({"_id": ObjectId("5994b2d48f244f71eb4d8569")})

// db.credentials.update(
//     {"_id": ObjectId("5994b2d48f244f71eb4d8569")},
//     {$set: {"passwordMD5HashHacked": true}}
// )
