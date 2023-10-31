package handlers

import (
	"casserole/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *BaseHandler) InternalReadHandler(c *fiber.Ctx) error {
	if h.NodeManager.Me().IsDead() {
		return c.SendStatus(503)
	}
	// Parse parameters
	courseId := c.Params("courseId")
	studentId := c.Params("studentId")

	data := internalRead(h.NodeManager, courseId, studentId)
	if data == nil {
		return c.SendStatus(404)
	}
	return c.JSON(data)
}

func internalRead(nm *utils.NodeManager, courseId, studentId string) *utils.Row {
	data, err := nm.DatabaseManager.GetRowByPartitionKey(courseId, studentId)
	if err != nil {
		return nil
	}
	return data
}
