package handlers

import (
	"casserole/utils"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

const BASE_READ_URL = "http://localhost:%d/read/%v"

func (h *BaseHandler) ReadHandler(c *fiber.Ctx) error {
	courseId := c.Params("courseId")

	/* get list of node ids to forward request to from CH */
	nodes := h.NodeManager.GetNodesForKey(courseId)

	noOfAck := 0
	reqsToForward := []utils.Request{}

	for _, node := range(nodes) {
		if node.Id == h.NodeManager.Id {
			// TODO: Read from self
			noOfAck++
			continue
		}

		reqsToForward = append(
			reqsToForward,
			utils.Request{
				NodeId: node.Id,
				Url: fmt.Sprintf(BASE_READ_URL, node.Port, courseId),
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

	//TODO: run read repair here

	if noOfAck >= h.NodeManager.Quorum {
		//return successful response with latest data
	}

	//return failed response status 500

	return nil
}
