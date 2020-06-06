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

// CreateUser is used by admin
func (h *Handler) CreateUser(c echo.Context) (err error) {
	return h.UserHandler.CreateUser(c)
}

// UpdateUser is a user information
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

func (h *Handler) ForgetPassword(c echo.Context) (err error) {
	return h.UserHandler.ForgetPassword(c)
}

func (h *Handler) ResetPassword(c echo.Context) (err error) {
	return h.UserHandler.ResetPassword(c)
}

func (h *Handler) ChangePassword(c echo.Context) (err error) {
	return h.UserHandler.ChangePassword(c)
}
