package handlers

import (
	"casserole/utils"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BaseHandler struct {
	NodeManager *utils.NodeManager
}

func NewHandler(nm *utils.NodeManager) *BaseHandler {
	return &BaseHandler{NodeManager: nm}
}

func (h *BaseHandler) DelayIfDead(c *fiber.Ctx) error {
	if !h.NodeManager.Me().IsDead() {
		return nil
	}

	done := make(chan bool)

	// Start the asynchronous sleep
	go func() {
		time.Sleep((h.NodeManager.GetConfig().Timeout + 1) * time.Second) // Sleep for 2 seconds
		done <- true
	}()

	// Non-blocking check if sleep is done
	select {
	case <-done:
		return errors.New("node is dead")
	}
}
