package handlers

import (
	"casserole/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h *BaseHandler) WriteHandler(c *fiber.Ctx) error {
	courseId := c.Params("courseId")

	/* get list of node ids to forward request to from CH */
	ids := []int{1, 2, 3} //TODO: remove after ch successors implementation

	nodes := make([]utils.Node, len(ids))
	for i, nodeId := range ids {
		node := h.NodeManager.ConfigManager.FindNodeById(nodeId)
		nodes[i] = *node
	}

	noOfAck := 0
	reqsToForward := []utils.Request{}

	//check to write to self
	for i, nodeId := range ids {

		//write to self
		if nodeId == h.NodeManager.Id {

			/* TODO: write to self */
			noOfAck++
			continue
		}

		reqsToForward = append(reqsToForward, utils.Request{
			NodeId: nodeId,
			Url:    fmt.Sprintf("http://localhost:%d/write/%d", nodes[i].Port, courseId)})
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
	if noOfAck == len(ids) {
		//return successful response with data
	}

	//some nodes respond
	//hinted handoff and successful response

	return nil
}
