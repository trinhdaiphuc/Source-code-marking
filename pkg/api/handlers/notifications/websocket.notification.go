package notifications

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
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
	Jwt            string `json:"jwt,omitempty"`
	NotificationID string `json:"notification_id,omitempty"`
	Notifications  string `json:"notifications,omitempty"`
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

	if h.WebsocketClients[ws] == "" {
		h.WebsocketClients[ws] = claims.Email
	}

	ctx := context.Background()
	h.Logger.Info("Connect to websocket user: ", claims.Email)

	notificationCollection := models.GetNotificationCollection(h.DB)

	opts := []*options.FindOptions{}
	opts = append(opts, options.Find().SetSort(bson.D{{"created_at", 1}}))
	opts = append(opts, options.Find().SetSkip(0))
	opts = append(opts, options.Find().SetLimit(5))
	opts = append(opts, options.Find().SetProjection(bson.D{{"content", 1}, {"exercise_id", 1}, {"is_read", 1}}))

	filter := bson.M{"user_id": claims.ID, "is_deleted": false}

	cursor, err := notificationCollection.Find(ctx, filter, opts...)
	if err != nil {
		h.Logger.Error("Error when find ", err)
		return
	}

	notificationArray := []models.Notification{}
	cursor.All(ctx, &notificationArray)
	data, _ := json.Marshal(notificationArray)

	firstMsg := &WebsocketMessage{
		Notifications: string(data),
	}

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
		// Read
		readMsg := &WebsocketMessage{}
		err = ws.ReadJSON(&readMsg)
		if err != nil {
			delete(h.WebsocketClients, ws)
			return err
		}

		writeMsg := &WebsocketMessage{}
		// Consume messages.
		for msgRedis := range ch {
			h.Logger.Debug(msgRedis.Channel, " ", msgRedis.Payload)
			// Write
			writeMsg.Notifications = msgRedis.Payload
			err = ws.WriteJSON(writeMsg)
			if err != nil {
				h.Logger.Error(err)
				delete(h.WebsocketClients, ws)
				return
			}
		}
	}
}
