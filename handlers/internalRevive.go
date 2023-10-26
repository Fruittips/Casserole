package handlers

import (
	"casserole/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const BASE_INTERNAL_REVIVE_URL = "http://localhost:%d/internal/revive/%v"

func (h *BaseHandler) InternalReviveHandler(c *fiber.Ctx) error {

	resp := InternalRevive(h.NodeManager)
	if resp.Error == nil && resp.StatusCode == http.StatusOK {
		return c.JSON(resp.Data)
	}
	return c.SendStatus(resp.StatusCode)
}

func InternalRevive(nm *utils.NodeManager) utils.Response {

	nm.Me().MakeAlive()

	if !nm.Me().IsDead() {
		return utils.Response{
			Error:      errors.New("isDead was not changed to false"),
			StatusCode: 500,
			NodeId:     nm.LocalId,
		}
	}

	// broadcast to all messages
	// find all the avail nodes in sysconfig

	for id, nodeData := range nm.Nodes {

		req := utils.Request{
			NodeId: id,
			Url:    fmt.Sprintf(BASE_INTERNAL_CHECKHH_URL, nodeData.Port, id),
		}
		res := nm.IntraSystemRequests([]utils.Request{req})
		for _, r := range res {
			if r.StatusCode == 500 {
				return utils.Response{
					StatusCode: 500,
					NodeId:     nm.LocalId,
				}
			}
		}
	}

	return utils.Response{
		StatusCode: 200,
		NodeId:     nm.LocalId,
	}
}
