package main

import (
	"casserole/handlers"
	"casserole/utils"
	"flag"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
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
	
	// Setup routes based on fstrings
	// These routes are located in the handler files
	read_endpoint_route  := fmt.Sprintf(utils.READ_ENDPOINT_FSTRING, ":courseId", ":studentId")
	write_endpoint_route := fmt.Sprintf(utils.WRITE_ENDPOINT_FSTRING, ":courseId", ":studentId")
	internal_read_endpoint_route  := fmt.Sprintf(utils.INTERNAL_READ_ENDPOINT_FSTRING, ":courseId", ":studentId")
	internal_write_endpoint_route := fmt.Sprintf(utils.INTERNAL_WRITE_ENDPOINT_FSTRING, ":courseId", ":studentId")
	internal_checkhh_endpoint_route := fmt.Sprintf(utils.INTERNAL_CHECKHH_ENDPOINT_FSTRING, ":nodeId")
	internal_kill_endpoint_route := utils.INTERNAL_KILL_ENDPOINT_FSTRING
	internal_revive_endpoint_route := utils.INTERNAL_REVIVE_ENDPOINT_FSTRING

	app := fiber.New()
	log.Printf("Node %v initialised on port %v.", nodeManager.LocalId, nodeManager.Me().Port)

	app.Use(logger.New())

	app.Get("/health-check", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World 👋!")
	})


	app.Get(read_endpoint_route, baseHandler.ReadHandler)
	
	app.Post(write_endpoint_route, baseHandler.WriteHandler)

	app.Get(internal_read_endpoint_route, baseHandler.InternalReadHandler)

	app.Post(internal_write_endpoint_route, baseHandler.InternalWriteHandler)

	app.Get(internal_checkhh_endpoint_route, baseHandler.InternalCheckHHHandler)

	app.Get(internal_revive_endpoint_route, baseHandler.InternalReviveHandler)

	app.Get(internal_kill_endpoint_route, baseHandler.InternalKillHandler)

	app.Listen(fmt.Sprintf(":%d", *port))
}
