package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"time"
)

func (h *Handler) HandleNotificationStream(c *gin.Context) {
	const op = "UserHandler.HandleNotificationStream"
	log := h.log.With(slog.String("op", op))

	// Set the response header to indicate SSE content type
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	userIDFromCtx, ok := c.Get("userID")
	if !ok {
		log.Warn("userID not found")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID := userIDFromCtx.(int64)

	ch := h.notificationService.AddNewConnection(userID)

	notify := c.Writer.CloseNotify()
	go func() {
		<-notify
		h.notificationService.RemoveConnection(userID)
		log.Info("Client disconnected")
	}()

	for {
		select {
		case <-notify:
			return
		case <-time.After(time.Second * 2):
			fmt.Fprintf(c.Writer, "data: %s\n\n", "lol")
			c.Writer.Flush()
		case notification := <-ch:
			log.Info("Notification received", slog.AnyValue(notification))
			fmt.Fprintf(c.Writer, "data: %v\n\n", notification)
			c.Writer.Flush()
		}
	}

	/*c.Stream(func(w io.Writer) bool {
		if notification := <-ch; ok {
			log.Info("Sending notification to client")
			c.SSEvent("notification", notification)
			return true
		}
		return false
	})*/

}
