package handlers

import (
	"casserole/utils"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const INTERNAL_REVIVE_ENDPOINT_FSTRING = "/internal/revive"

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
	// Internal Check Hinted Handoffs URL: Port
	internal_checkhh_url := "http://localhost:%d" + INTERNAL_CHECKHH_ENDPOINT_FSTRING

	for id, nodeData := range nm.Nodes {

		req := utils.Request{
			NodeId: id,
			Url:    fmt.Sprintf(internal_checkhh_url, nodeData.Port),
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
