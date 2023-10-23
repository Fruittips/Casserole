# How to run
```
go run . -port=3001
```
# How to create an endpoint
Every request needs to be able to use the NodeManager, we do this via dependency injection. 
The dependency, NodeManager, is injected into the handler via the BaseHandler struct.
Therefore, every request handler needs to be a method of the BaseHandler struct.
```go
package handlers

import "github.com/gofiber/fiber/v2"

func (h *BaseHandler) WriteHandler(c *fiber.Ctx) error {
  // how to get access to the node manager
  node := h.NodeManager

  return nil
}
```
```go
package handlers

import "casserole/utils"

type BaseHandler struct {
	NodeManager *utils.NodeManager
}

func NewHandler(nm *utils.NodeManager) *BaseHandler {
	return &BaseHandler{NodeManager: nm}
}

```

The NodeManager and BaseHandler are instantiated at the start (main.go)
```go
package main

import (
	"casserole/handlers"
	"casserole/utils"
	"flag"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var port = flag.Int("port", 8080, "port to listen on")
var nodeManager *utils.NodeManager

func main() {
	flag.Parse()
	nodeManager = utils.NewNodeManager(*port)
	baseHandler := handlers.NewHandler(nodeManager)

	app := fiber.New()

	app.Get("/health-check", func(ctx *fiber.Ctx) error {
		return ctx.SendString("Hello, World 👋!")
	})

	app.Get("/write", baseHandler.WriteHandler)

	app.Listen(fmt.Sprintf(":%d", port))
}

```

# Managers
### Node Manager
```go
node := utils.NewNodeManager(3001)

// append to hinted handoffs
node.HintedHandoffManager.Append(1, utils.AtomicDbMessage{Data: []string{"hello", "world"}, Timestamp: 123})

// append to db
node.DatabaseManager.AppendRow(utils.Row{Data: []string{"hello", "asdads"}, Timestamp: 123})

// read from db
fmt.Println(node.DatabaseManager.Data)

// get row by id
data, err := node.DatabaseManager.GetRowById(1)
fmt.Println(data)

// read from hinted handoffs
fmt.Println(node.HintedHandoffManager.Data)

// get config data
fmt.Println(node.ConfigManager.Data)

// get my node data
fmt.Println(node.Me())
```




