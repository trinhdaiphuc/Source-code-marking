package handlers

import "github.com/labstack/echo/v4"

/////////////////////////////////////////////////////////
// TODO Notification
/////////////////////////////////////////////////////////

// WebsocketNotification create a websocket notification connection
func (h *Handler) WebsocketNotification(c echo.Context) (err error) {
	return h.Notification.WebsocketNotification(c)
}

func (h *Handler) GetAllNotifications(c echo.Context) (err error) {
	return h.Notification.GetAllNotifications(c)
}

func (h *Handler) MarkReadNotification(c echo.Context) (err error) {
	return h.Notification.MarkReadNotification(c)
}
