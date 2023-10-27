package handlers

import (
	"casserole/utils"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (h *BaseHandler) ReadHandler(c *fiber.Ctx) error {
	// Internal read URL: Port, CourseId, StudentId
	internal_read_url := "http://localhost:%d" + INTERNAL_READ_ENDPOINT_FSTRING

	courseId := c.Params("courseId")
	studentId := c.Params("studentId")

	/* get list of node ids to forward request to from CH */
	nodes := h.NodeManager.GetNodesForKey(courseId)

	noOfAck := 0
	reqsToForward := []utils.Request{}

	responses := []utils.Response{}

	for _, node := range nodes {
		log.Printf("Node %v: READ(%v, %v) from node %v", h.NodeManager.LocalId, courseId, studentId, node.Id)

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
				Url:    fmt.Sprintf(internal_read_url, node.Port, courseId, studentId),
			},
		)
	}

	latestRecord := utils.Row{}
	responses = append(responses, h.IntraSystemRequests(reqsToForward)...)
	for _, res := range responses {
		if res.Error != nil {
			continue
		}

		if res.Data.CreatedAt > latestRecord.CreatedAt || latestRecord == (utils.Row{}) {
			latestRecord = *res.Data
		}
		noOfAck++
	}

	//TODO: run read repair here [untested]
	rrm := utils.NewReadRepairsManager(h.NodeManager, courseId, studentId, responses)
	rrm.PerformReadRepair(responses)

	if noOfAck >= h.NodeManager.Quorum {
		//return successful response with latest data
		return c.JSON(latestRecord)
	}

	//return failed response status 500

	return c.SendStatus(500)
}
