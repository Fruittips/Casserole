package handlers

import (
	"casserole/utils"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func (h *BaseHandler) WriteHandler(c *fiber.Ctx) error {
	// Internal write URL: Port, CourseId, StudentId
	internal_write_url := "http://localhost:%d" + INTERNAL_WRITE_ENDPOINT_FSTRING
	
	courseId := c.Params("courseId")

	newStudent := utils.Row{}
	err := c.BodyParser(&newStudent)
	if err != nil {
		return err
	}

	/* get list of node ids to forward request to from CH */
	nodes := h.NodeManager.GetNodesForKey(courseId)

	noOfAck := 0
	reqsToForward := []utils.Request{}

	for _, node := range nodes {
		log.Printf("Node %v: WRITE(%v, %v) to node %v with data: %v", h.NodeManager.LocalId, courseId, newStudent.StudentId, node.Id, newStudent)
		if node.Id == h.NodeManager.LocalId {
			err := h.NodeManager.DatabaseManager.AppendRow(courseId, newStudent)
			if err != nil {
				noOfAck++
			}
			continue
		}

		reqsToForward = append(
			reqsToForward,
			utils.Request{
				NodeId:  node.Id,
				Url:     fmt.Sprintf(internal_write_url, node.Port, courseId, newStudent.StudentId),
				Payload: &newStudent,
			},
		)
	}

	responses := h.NodeManager.IntraSystemRequests(reqsToForward)
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
