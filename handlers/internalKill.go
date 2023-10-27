package handlers

import (
	"casserole/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func (h *BaseHandler) InternalKillHandler(c *fiber.Ctx) error {
	err := internalKill(h.NodeManager)
	if err != nil {
		return c.SendStatus(500)
	}
	return c.SendStatus(200)
}

func internalKill(nm *utils.NodeManager) error {
	nm.Me().MakeDead()

	if !nm.Me().IsDead() {
		return errors.New("isDead not changed to true")
	}
	return nil
}
