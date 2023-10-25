package main

import (
	"casserole/handlers"
	"casserole/utils"
	"flag"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var nodeManager *utils.NodeManager

func main() {
	port := flag.Int("port", -1, "port to listen on")
	flag.Parse()
	if *port < 0 {
		panic("port is required")
	}
	nodeManager = utils.NewNodeManager(*port)
	baseHandler := handlers.NewHandler(nodeManager)

	app := fiber.New()

	app.Get("/health-check", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World 👋!")
	})

	app.Get("/write/:courseId/student/:studentId", baseHandler.WriteHandler)

	app.Get("/read/:courseId/student/:studentId", baseHandler.ReadHandler)

	app.Get("/internal/read/courses/:courseId/student/:studentId", baseHandler.InternalReadHandler)

	app.Get("/internal/write/courses/:courseId/student/:studentId", baseHandler.InternalWriteHandler)

	app.Listen(fmt.Sprintf(":%d", *port))
}
