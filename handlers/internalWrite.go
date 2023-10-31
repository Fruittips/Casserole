package handlers

import (
	"casserole/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *BaseHandler) InternalWriteHandler(c *fiber.Ctx) error {
	if h.NodeManager.Me().IsDead() {
		return c.SendStatus(503)
	}
	// Parse parameters
	courseId := c.Params("courseId")

	// Parse body -- we're given a Row in JSON
	newRow := new(utils.Row)
	if err := c.BodyParser(newRow); err != nil {
		return err
	}

	// Perform local write
	err := internalWrite(h.NodeManager, courseId, *newRow)
	if err != nil {
		return c.SendStatus(500)
	}
	return c.SendStatus(200)
}

func internalWrite(nm *utils.NodeManager, courseId string, newData utils.Row) error {
	err := nm.DatabaseManager.AppendRow(courseId, newData)
	return err
}
