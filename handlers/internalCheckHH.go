package handlers

import (
	"casserole/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *BaseHandler) InternalCheckHHHandler(c *fiber.Ctx) error {
	resp := InternalCheckHH(h.NodeManager)
	if resp.Error == nil && resp.StatusCode == http.StatusOK {
		return c.JSON(resp.Data)
	}
	return c.SendStatus(resp.StatusCode)
}

func InternalCheckHH(nm *utils.NodeManager) utils.Response {
	// for logging
	for id, outerrow := range nm.HintedHandoffManager.Data.Rows {
		if string(nm.Me().Port) == id {
			for i, adbm := range outerrow {
				InternalWrite(nm, id, adbm.Data[i])
			}
		}
	}

	return utils.Response{
		StatusCode: 200,
		NodeId:     nm.LocalId,
	}
}
