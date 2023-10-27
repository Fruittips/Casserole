package handlers

import (
	"casserole/utils"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type internalReadResponse struct {
	srcId utils.NodeId
	data  *utils.Row
}

func (h *BaseHandler) ReadHandler(c *fiber.Ctx) error {
	courseId := c.Params("courseId")
	studentId := c.Params("studentId")

	// Get list of nodes from CHT
	nodes := h.NodeManager.GetNodesForKey(courseId)
	var reqWg sync.WaitGroup

	responses := make(chan internalReadResponse, len(nodes))
	for _, node := range nodes {
		log.Printf("Node %v: READ(%v, %v) from node %v", h.NodeManager.LocalId, courseId, studentId, node.Id)
		// Query self if self is one of the nodes
		if node.Id == h.NodeManager.LocalId {
			response := internalReadResponse{node.Id, internalRead(h.NodeManager, courseId, studentId)}
			responses <- response
			continue
		}

		// Otherwise, shoot a concurrent internal read
		reqWg.Add(1)
		go func(n *utils.Node) {
			defer reqWg.Done()
			row, err := h.NodeManager.SendInternalRead(*n, courseId, studentId)
			if err != nil {
				log.Printf("Node %v READ from node %v Error: %v", h.NodeManager.LocalId, n.Id, err)
				responses <- internalReadResponse{n.Id, nil}
				return
			}
			responses <- internalReadResponse{n.Id, row}
		}(node)
	}

	reqWg.Wait()
	close(responses)

	// Get latest record from buffered channel
	ackCount := 0
	var latestRecord *utils.Row
	responses_ls := make([]internalReadResponse, 0)
	for res := range responses {
		data := res.data
		responses_ls = append(responses_ls, res)
		if data == nil { // ignore if it's an empty response (i.e. no response)
			continue
		}

		// Otherwise, it's a legitimate response
		ackCount++
		if latestRecord == nil {
			latestRecord = data
			continue
		}

		// Only change latest if it's actually later
		if data.CreatedAt > latestRecord.CreatedAt {
			latestRecord = data
		}
	}

	// Identify nodes with outdated data, get them to read repair
	//TODO:
	// Ryan: Feels that ReadRepairsManager is unnecessary, no point creating a new manager for an ephemeral thing
	//       Instead, just use a separate fn in here to handle any potential discrepancies with the latest and continue on our merry way
	go readRepair(h.NodeManager, courseId, *latestRecord, responses_ls)
	//TODO
	// Alternative to above:
	//TODO: run read repair here [untested]
	//rrm := utils.NewReadRepairsManager(h.NodeManager, courseId, studentId, responses)
	//rrm.PerformReadRepair(responses)

	// Only return successful response if with quorum
	if ackCount >= h.NodeManager.Quorum {
		return c.JSON(latestRecord)
	}
	return c.SendStatus(500)
}

func readRepair(nm *utils.NodeManager, courseId string, latestRecord utils.Row, internalReadResponses []internalReadResponse) {
	//TODO: do in goroutines and wg
	for _, res := range internalReadResponses {
		if !latestRecord.Equal(*res.data) {
			outdatedNode, err := nm.GetNodeById(utils.NodeId(res.srcId))

			if err != nil {
				log.Printf("Node %v READ REPAIR node %v Error: %v", nm.LocalId, res.srcId, err)
				continue
			}
			log.Printf("Node %v READ REPAIR node %v.", nm.LocalId, outdatedNode.Id)
			err = nm.SendInternalWrite(*outdatedNode, courseId, latestRecord)

			if err != nil {
				log.Printf("Node %v READ REPAIR node %v Error: %v", nm.LocalId, res.srcId, err)
				continue
			}
		}
	}
}
