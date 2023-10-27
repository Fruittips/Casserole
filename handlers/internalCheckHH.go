package handlers

import (
	"casserole/utils"
	"github.com/gofiber/fiber/v2"
)

// When called, returns a row for the newly revived node
func (h *BaseHandler) InternalCheckHHHandler(c *fiber.Ctx) error {
	nodeIdToCheck := c.Params("nodeId")

	err := internalCheckHH(h.NodeManager, utils.NodeId(nodeIdToCheck))
	if err != nil {
		return c.SendStatus(404) // TODO: Two possible errors -- one is error in the function itself, another is no hinted handoffs found for the node.
	}
	return c.SendStatus(200) // Since we're sending an internal write separately.
}

func internalCheckHH(nm *utils.NodeManager, nodeIdToCheck utils.NodeId) error {
	for id, outerrow := range nm.HintedHandoffManager.Data.Rows {
		if nodeIdToCheck == utils.NodeId(id) {
			// Current row refers to the node ID to be checked
			for i, adbm := range outerrow {
				targetNode, err := nm.GetNodeById(nodeIdToCheck)
				if err != nil {
					return err
				}
				courseId := "TODO"
				data := adbm.Data[i]

				// Send an internal write to that node
				err = nm.SendInternalWrite(
					*targetNode,
					courseId,
					data,
				)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
	
}
