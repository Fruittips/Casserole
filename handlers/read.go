package handlers

import (
	"casserole/utils"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
)

func (h *BaseHandler) ReadHandler(c *fiber.Ctx) error {
	courseId := c.Params("courseId")
	studentId := c.Params("studentId")

	// Get list of nodes from CHT
	nodes := h.NodeManager.GetNodesForKey(courseId)
	var reqWg sync.WaitGroup

	responses := make(chan *utils.Row, len(nodes))
	for _, node := range nodes {
		log.Printf("Node %v: READ(%v, %v) from node %v", h.NodeManager.LocalId, courseId, studentId, node.Id)

		// Query self if self is one of the nodes
		if node.Id == h.NodeManager.LocalId {
			responses <- internalRead(h.NodeManager, courseId, studentId)
			continue
		}

		// Otherwise, shoot a concurrent internal read
		reqWg.Add(1)
		go func(n *utils.Node) {
			defer reqWg.Done()
			row, err := h.NodeManager.SendInternalRead(*n, courseId, studentId)
			if err != nil {
				log.Printf("Node %v READ from node %v Error: %v", h.NodeManager.LocalId, n.Id, err)
				return
			}
			responses <- row
		}(node)
	}

	reqWg.Wait()
	close(responses)

	// Get latest record from buffered channel
	ackCount := 0
	var latestRecord *utils.Row
	for res := range responses {
		ackCount++
		if latestRecord == nil {
			latestRecord = res
			continue
		}

		// Only change latest if it's actually later
		if res.CreatedAt > latestRecord.CreatedAt {
			latestRecord = res
		}
	}

	// TODO: Read Repair

	// Only return successful response if with quorum
	if ackCount >= h.NodeManager.Quorum {
		return c.JSON(latestRecord)
	}
	return c.SendStatus(500)
}
