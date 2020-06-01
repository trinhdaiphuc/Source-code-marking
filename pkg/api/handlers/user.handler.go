package handlers

import "github.com/labstack/echo/v4"

/////////////////////////////////////////////////////////
// TODO User handler
/////////////////////////////////////////////////////////
// ShowAccount godoc
// @Summary Show a account
// @Description get string by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {object} model.Account
// @Header 200 {string} Token "qwerty"
// @Failure 400 {object} httputil.HTTPError
// @Failure 404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Router /accounts/{id} [get]

// Signup an account
func (h *Handler) Signup(c echo.Context) (err error) {
	return h.UserHandler.Signup(c)
}

// Signin to an account
func (h *Handler) Signin(c echo.Context) (err error) {
	return h.UserHandler.Signin(c)
}

// Profile is a user information
func (h *Handler) Profile(c echo.Context) (err error) {
	return h.UserHandler.Profile(c)
}

// GetUser is a user information
func (h *Handler) GetUser(c echo.Context) (err error) {
	return h.UserHandler.GetUser(c)
}

// GetAllUsers is a user information
func (h *Handler) GetAllUsers(c echo.Context) (err error) {
	return h.UserHandler.GetAllUsers(c)
}

// UpdateUsers is a user information
func (h *Handler) UpdateUser(c echo.Context) (err error) {
	return h.UserHandler.UpdateUser(c)
}

// DeleteUser is a user information
func (h *Handler) DeleteUser(c echo.Context) (err error) {
	return h.UserHandler.DeleteUser(c)
}

// ListClasses is a user information
func (h *Handler) ListClasses(c echo.Context) (err error) {
	return h.UserHandler.ListClasses(c)
}

func (h *Handler) ValidateUser(c echo.Context) (err error) {
	return h.UserHandler.ValidateUser(c)
}
