package handlers

import (
	"casserole/utils"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

const BASE_WRITE_URL = "http://localhost:%d/write/%v"

func (h *BaseHandler) WriteHandler(c *fiber.Ctx) error {
	courseId := c.Params("courseId")

	/* get list of node ids to forward request to from CH */
	nodes := h.NodeManager.GetNodesForKey(courseId)
	// for logging
	for _, node := range(nodes) {
		log.Printf("Writing %v to N%d", courseId, node.Id)
	}

	noOfAck := 0
	reqsToForward := []utils.Request{}

	for _, node := range(nodes) {
		if node.Id == h.NodeManager.LocalId {
			// TODO: Write from self
			noOfAck++
			continue
		}

		reqsToForward = append(
			reqsToForward,
			utils.Request{
				NodeId: int(node.Id),
				Url: fmt.Sprintf(BASE_WRITE_URL, node.Port, courseId),
			},
		)
	}

	responses := h.NodeManager.ForwardGetRequests(reqsToForward)
	for _, res := range responses {
		if res.Error != nil {
			continue
		}

		//TODO: get the last written value
		noOfAck++
	}

	// if failed to hit QUORUM
	if noOfAck < h.NodeManager.Quorum {
		/* TODO: write to hinted handoff */
		//return with error response 500
	}

	/* if hit quorum
	1. all nodes respond
	2. some nodes respond
	*/

	//if all nodes respond
	if noOfAck == len(nodes) {
		//return successful response with data
	}

	//some nodes respond
	//hinted handoff and successful response

	return c.SendStatus(500)
}
