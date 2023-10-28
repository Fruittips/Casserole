package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func (h *BaseHandler) SseEndpointHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		for {
			nodeData := h.NodeManager.Me().String()
			dbData := h.NodeManager.DatabaseManager.Data.String()
			hhData := fmt.Sprintf(h.NodeManager.HintedHandoffManager.Data.String())

			data := map[string]string{
				"node": nodeData,
				"db":   dbData,
				"hh":   hhData,
			}

			jsonData, err := json.Marshal(data)
			if err != nil {
				// Handle error
				fmt.Println("Error marshalling data for SSE endpoint: ", err)
				return
			}

			ssePayload := fmt.Sprintf("data: %s\n\n", jsonData)

			_, err = w.WriteString(ssePayload)
			if err != nil {
				// Handle error
				fmt.Println("Error writing SSE payload: ", err)
				return
			}

			err = w.Flush()
			if err != nil {
				// An error occurred, the client might have disconnected
				fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)
				return
			}

			time.Sleep(2 * time.Second)
		}
	}))

	return nil
}

func generateSSEMessage(data fmt.Stringer) []byte {
	return []byte("data: " + data.String() + "\n\n")
}
