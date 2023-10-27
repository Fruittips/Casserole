package handlers

import (
	"casserole/utils"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func (h *BaseHandler) WriteHandler(c *fiber.Ctx) error {
	courseId := c.Params("courseId")

	newStudent := utils.Row{}
	err := c.BodyParser(&newStudent)
	if err != nil {
		return err
	}

	// Get list of nodes from CHT
	nodes := h.NodeManager.GetNodesForKey(courseId)
	var reqWg sync.WaitGroup

	responses := make(chan bool, len(nodes))
	for _, node := range nodes {
		log.Printf("Node %v: WRITE(%v, %v) to node %v with data: %v", h.NodeManager.LocalId, courseId, newStudent.StudentId, node.Id, newStudent)

		// Query self if self is one of hte nodes
		if node.Id == h.NodeManager.LocalId {
			err := internalWrite(h.NodeManager, courseId, newStudent)
			if err != nil {
				log.Printf("Node %v WRITE to node %v Error: %v", h.NodeManager.LocalId, node.Id, err)
			} else {
				responses <- true
			}
			continue
		}

		// Otherwise, shoot a concurrent internal write
		reqWg.Add(1)
		go func(n *utils.Node) {
			defer reqWg.Done()
			err := h.NodeManager.SendInternalWrite(*n, courseId, newStudent)
			if err != nil {
				log.Printf("Node %v WRITE to node %v Error: %v", h.NodeManager.LocalId, n.Id, err)
				return
			}
			responses <- true
		}(node)
	}

	reqWg.Wait()
	close(responses)

	// Check number of acks
	ackCount := len(responses)

	// If failed to hit QUORUM
	if ackCount < h.NodeManager.Quorum {
		// TODO: Write to hinted handoff
		return c.SendStatus(500)
	}

	// If hit QUORUM:
	// Either all nodes responded, or some nodes responded
	if ackCount == len(nodes) {
		return c.SendStatus(200)
	} else {
		// TODO: Hinted Handoffs

		return c.SendStatus(200)
	}
}
