package handlers

import "github.com/labstack/echo/v4"

/////////////////////////////////////////////////////////
// TODO Notification
/////////////////////////////////////////////////////////

// WebSocketNotification create a websocket notification connection
func (h *Handler) WebsocketNotification(c echo.Context) (err error) {
	return h.Notification.WebsocketNotification(c)
}
