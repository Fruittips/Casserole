package handlers

import (
	"casserole/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type AtomicDbMessage struct {
	Data      []string `json: "Data"`
	Timestamp int64    `json: "Timestamp"`
}

type Data struct {
	TableName string                     `json: "TableName"`
	Columns   []string                   `json: "Columns"`
	Row       map[string]AtomicDbMessage `json: "Row"`
}

const BASE_INTERNAL_READ_URL = "http://localhost:%d/internal/read/%v"

func (h *BaseHandler) InternalReadHandler(c *fiber.Ctx) error {
	r := new(utils.Request)
	if err := c.BodyParser(r); err != nil {
		return err
	}

	resp := InternalRead(h.NodeManager, r.CourseId, r.StudentId)
	if resp.Error == nil && resp.StatusCode == http.StatusOK {
		return c.JSON(resp.Data)
	}
	return c.SendStatus(resp.StatusCode)
}

func InternalRead(nm *utils.NodeManager, courseId string, studentId string) utils.Response {
	data, err := nm.DatabaseManager.GetRowByPartitionKey(courseId, studentId)
	if err != nil {
		return utils.Response{
			Error:      err,
			StatusCode: 500,
			NodeId:     nm.LocalId,
		}
	}

	return utils.Response{
		Data:       data,
		StatusCode: 200,
		NodeId:     nm.LocalId,
	}
}
