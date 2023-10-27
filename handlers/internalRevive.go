package handlers

import (
	"casserole/utils"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func (h *BaseHandler) InternalReviveHandler(c *fiber.Ctx) error {

	err := internalRevive(h.NodeManager)
	if err != nil {
		return c.SendStatus(500)
	}
	return c.SendStatus(200)
}

func internalRevive(nm *utils.NodeManager) error {
	log.Printf("Node %v revived.", nm.LocalId)
	
	nm.Me().MakeAlive()

	if nm.Me().IsDead() {
		return errors.New("isDead not changed to false")
	}

	// Request all other nodes for hintedhandoffs
	var reqWg sync.WaitGroup
	responses := make(chan error, len(nm.Nodes) - 1) // Ignore self
	for _, node := range nm.Nodes {
		if node.Id == nm.LocalId {
			continue
		}

		reqWg.Add(1)
		go func(n *utils.Node) {
			defer reqWg.Done()
			err := nm.RequestForHHs(*n)
			if err != nil {
				responses <- errors.New(fmt.Sprintf("Node %v REQUEST HH from node %v Error: %v", nm.LocalId, n.Id, err))
				return
			}
			responses <- nil
		}(node)
	}

	reqWg.Wait()
	close(responses)

	// TODO: Any other or more specific error condition?
	for err := range responses {
		if err != nil {
			return err
		}
	}
	
	return nil
}
