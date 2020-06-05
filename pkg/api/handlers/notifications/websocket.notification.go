package notifications

import (
	"context"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
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
	Jwt     string `json:"jwt"`
	Message string `json:"message"`
}

func (h *NotificationHandler) WebsocketNotification(c echo.Context) (err error) {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	msg := &WebsocketMessage{}
	_, message, err := ws.ReadMessage()
	if err != nil {
		return err
	}
	jwtString := string(message[:])
	h.Logger.Debug("JWT ", jwtString)
	if jwtString == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Missing token",
		}
	}

	// Initialize a new instance of `Claims`
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(jwtString, claims,
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
	h.Logger.Info("Connect to websocket user: ", claims.Email)

	pubsub := h.RedisClient.Subscribe(ctx, claims.Email)
	// Wait for confirmation that subscription is created before publishing anything.
	_, err = pubsub.Receive(ctx)
	if err != nil {
		h.Logger.Error("Error when received redis ", err)
	}

	ch := pubsub.Channel()
	for {
		// Consume messages.
		for msgRedis := range ch {
			h.Logger.Debug(msgRedis.Channel, " ", msgRedis.Payload)
			// Read
			msg.Message = msgRedis.Payload
			err = ws.WriteJSON(msg)
			if err != nil {
				h.Logger.Error(err)
				delete(h.WebsocketClients, ws)
				return
			}
		}
	}
}
