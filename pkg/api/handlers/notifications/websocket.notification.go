package notifications

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		if os.Getenv("ENV") == "production" {
			if r.Header.Get("Origin") == os.Getenv("FRONT_END_SERVER_HOST") {
				return true
			}
			return false
		}
		return true
	},
}

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

// Define our message object
type WebsocketMessage struct {
	Jwt           string `json:"jwt,omitempty"`
	Notifications string `json:"notifications,omitempty"`
}

func (h *NotificationHandler) WebsocketNotification(c echo.Context) (err error) {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	msg := &WebsocketMessage{}
	err = ws.ReadJSON(&msg)
	if err != nil {
		return err
	}

	if msg.Jwt == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Missing token",
		}
	}

	// Initialize a new instance of `Claims`
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(msg.Jwt, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
	if err != nil {
		if err.Error() == jwt.ErrSignatureInvalid.Error() {
			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "Error token",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Error token",
			Internal: err,
		}
	}
	if !tkn.Valid {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Error token",
			Internal: err,
		}
	}

	h.WebsocketClients[ws] = claims.Email

	ctx := context.Background()
	h.Logger.Debug("Connect to websocket user: ", claims.Email)

	filter := bson.M{"user_id": claims.ID, "is_deleted": false}

	listParam := models.ListQueryParam{
		PageSize:  5,
		PageToken: 1,
		OrderBy:   "created_at",
		OrderType: internal.DESC.String(),
	}

	listNotification, err := models.ListAllNotifications(h.DB, filter, listParam)
	notificationCollection := models.GetNotificationCollection(h.DB)
	filter["is_read"] = false
	totalUnread, err := notificationCollection.CountDocuments(context.TODO(), filter)
	listNotificationWebsocket := models.ListNotificationWebsocket{
		Notifications: listNotification.Notifications,
		TotalUnread:   totalUnread,
	}
	data, _ := json.Marshal(listNotificationWebsocket)

	firstMsg := &WebsocketMessage{
		Notifications: string(data),
	}

	h.Logger.Debug("Message send ", firstMsg)

	err = ws.WriteJSON(firstMsg)
	if err != nil {
		h.Logger.Error(err)
		delete(h.WebsocketClients, ws)
		return
	}

	pubsub := h.RedisClient.Subscribe(ctx, claims.Email)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err = pubsub.Receive(ctx)
	if err != nil {
		h.Logger.Error("Error when received redis ", err)
	}

	ch := pubsub.Channel()
	defer pubsub.Close()

	for {
		writeMsg := &WebsocketMessage{}
		// Consume messages.
		for msgRedis := range ch {
			h.Logger.Debug("Message in channel ", msgRedis.Channel, " ", msgRedis.Payload)
			// Write
			writeMsg.Notifications = msgRedis.Payload
			err = ws.WriteJSON(writeMsg)
			if err != nil {
				h.Logger.Error("Error write socket fail ", err)
				delete(h.WebsocketClients, ws)
				break
			}
		}
	}
}
