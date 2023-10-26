package handlers

import (
	"casserole/utils"
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

const BASE_INTERNAL_KILL_URL = "http://localhost:%d/internal/kill/%v"

func (h *BaseHandler) InternalKillHandler(c *fiber.Ctx) error {

	resp := InternalKill(h.NodeManager)
	if resp.Error == nil && resp.StatusCode == http.StatusOK {
		return c.JSON(resp.Data)
	}
	return c.SendStatus(resp.StatusCode)
}

func InternalKill(nm *utils.NodeManager) utils.Response {

	nm.Me().MakeDead()

	if nm.Me().IsDead() {
		return utils.Response{
			Error:      errors.New("isDead was not changed to true"),
			StatusCode: 200,
			NodeId:     nm.LocalId,
		}
	}

	return utils.Response{
		StatusCode: 500,
		NodeId:     nm.LocalId,
	}

}
