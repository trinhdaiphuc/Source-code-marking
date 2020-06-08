package users

import (
	"context"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func sendResetPasswordEmail(user models.User, jwtKey string, logger *internal.AppLog) {
	token, _ := createTokenWithUser(user, jwtKey, 24)
	validationLink := os.Getenv("FRONT_END_SERVER_HOST") + "/password/" + token
	logger.Info("Validation link ", validationLink)
	content := "Please click this link to reset your password " + validationLink
	subject := "Welcome to Source code marking"
	id, err := internal.SendMail(os.Getenv("EMAIL_USERNAME"), user.Email, subject, content)
	if err != nil {
		logger.Error("Error when send mail ", err)
	} else {
		logger.Info("Send mail success with id ", id)
	}
}

func (h *UserHandler) ForgetPassword(c echo.Context) (err error) {
	email := c.QueryParam("email")
	h.Logger.Debug("Email ", email)
	userCollection := models.GetUserCollection(h.DB)
	ctx := context.Background()
	user := &models.User{}
	result := userCollection.FindOne(ctx, bson.M{"email": email})
	if err = result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found your email",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[ForgetPasword] Internal server error",
			Internal: err,
		}
	}

	if user.IsDeleted {
		return &echo.HTTPError{
			Code:    http.StatusGone,
			Message: "User has been deleted.",
		}
	}

	if user.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Wrong service")
	}

	go sendResetPasswordEmail(*user, h.JWTKey, h.Logger)
	return c.NoContent(http.StatusNoContent)
}
