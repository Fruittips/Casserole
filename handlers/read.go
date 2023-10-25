package handlers

import (
	"casserole/utils"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

const BASE_READ_URL = "http://localhost:%d/read/%v"

func (h *BaseHandler) ReadHandler(c *fiber.Ctx) error {
	courseId := c.Params("courseId")
	studentId := c.Params("studentId")

	/* get list of node ids to forward request to from CH */
	nodes := h.NodeManager.GetNodesForKey(courseId)
	for _, node := range nodes {
		log.Printf("Reading %v from node %v", courseId, node.Id)
	}

	noOfAck := 0
	reqsToForward := []utils.Request{}

	responses := []utils.Response{}

	for _, node := range nodes {
		if node.Id == h.NodeManager.LocalId {
			// Read from self
			r := InternalRead(h.NodeManager, courseId, studentId)
			responses = append(responses, r)
			noOfAck++
			continue
		}

		reqsToForward = append(
			reqsToForward,
			utils.Request{
				NodeId: node.Id,
				Url:    fmt.Sprintf(BASE_INTERNAL_READ_URL, node.Port, courseId),
			},
		)
	}

	latestRecord := utils.Row{}
	responses = append(responses, h.NodeManager.IntraSystemRequests(reqsToForward)...)
	for _, res := range responses {
		if res.Error != nil {
			continue
		}

		if res.Data.CreatedAt > latestRecord.CreatedAt || latestRecord == (utils.Row{}) {
			latestRecord = *res.Data
		}
		noOfAck++
	}

	//TODO: run read repair here

	if noOfAck >= h.NodeManager.Quorum {
		//return successful response with latest data
		return c.JSON(latestRecord)
	}

	//return failed response status 500

	return c.SendStatus(500)
}
