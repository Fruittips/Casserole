package handlers

import (
	"casserole/utils"
	"log"

	"github.com/gofiber/fiber/v2"
)

// When called, returns a row for the newly revived node
func (h *BaseHandler) InternalCheckHHHandler(c *fiber.Ctx) error {
	if h.NodeManager.Me().IsDead() {
		return c.SendStatus(503)
	}
	nodeIdToCheck := c.Params("nodeId")

	err := internalCheckHH(h.NodeManager, utils.NodeId(nodeIdToCheck))
	if err != nil {
		return c.SendStatus(404) // TODO: Two possible errors -- one is error in the function itself, another is no hinted handoffs found for the node.
	}
	return c.SendStatus(200) // Since we're sending an internal write separately.
}

func internalCheckHH(nm *utils.NodeManager, nodeIdToCheck utils.NodeId) error {
	nm.HintedHandoffManager.Mux.Lock()
	defer nm.HintedHandoffManager.Mux.Unlock()
	for id, outerrow := range nm.HintedHandoffManager.Data.Rows {
		if nodeIdToCheck == utils.NodeId(id) {
			// Current row refers to the node ID to be checked

			for len(outerrow) > 0 {
				adbm := outerrow[0]

				targetNode, err := nm.GetNodeById(nodeIdToCheck)
				if err != nil {
					return err
				}
				courseId := adbm.CourseId
				data := adbm.Data

				// Send an internal write to that node
				err = nm.SendInternalWrite(
					*targetNode,
					courseId,
					data,
				)
				if err != nil {
					log.Printf("Error in internal write for node %v", id)
					break
				}

				// pop
				outerrow = outerrow[1:]
			}
			nm.HintedHandoffManager.Data.Rows[id] = outerrow
		}
	}
	nm.HintedHandoffManager.OverwriteWithMem()
	return nil

}
