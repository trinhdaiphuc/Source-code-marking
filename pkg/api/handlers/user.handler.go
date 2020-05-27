package handlers

import "github.com/labstack/echo/v4"

/////////////////////////////////////////////////////////
// TODO User handler
/////////////////////////////////////////////////////////

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
