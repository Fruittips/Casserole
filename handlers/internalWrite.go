package handlers

import (
	"casserole/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const BASE_INTERNAL_WRITE_URL = "http://localhost:%d/internal/write/%v"

func (h *BaseHandler) InternalWriteHandler(c *fiber.Ctx, newData utils.Row) error {
	// failure response

	r := new(utils.Request)
	if err := c.BodyParser(r); err != nil {
		return err
	}

	resp := InternalWrite(h.NodeManager, r.CourseId, *r.Payload)
	if resp.Error == nil && resp.StatusCode == http.StatusOK {
		return c.JSON(resp.Data)
	}
	return c.SendStatus(resp.StatusCode)
}

func InternalWrite(nm *utils.NodeManager, partitionKey string, newData utils.Row) utils.Response {

	err := nm.DatabaseManager.AppendRow(partitionKey, newData)
	if err != nil {
		return utils.Response{
			Error:      err,
			StatusCode: 500,
			NodeId:     nm.LocalId,
		}
	}

	return utils.Response{
		StatusCode: 200,
		NodeId:     nm.LocalId,
	}
}
