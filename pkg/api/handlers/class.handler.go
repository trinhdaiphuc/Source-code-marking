package handlers

import "github.com/labstack/echo/v4"

/////////////////////////////////////////////////////////
// TODO Class handler
/////////////////////////////////////////////////////////

// CreateClass create a Class
func (h *Handler) CreateClass(c echo.Context) (err error) {
	return h.ClassHandler.CreateClass(c)
}

// UpdateClass update a Class
func (h *Handler) UpdateClass(c echo.Context) (err error) {
	return h.ClassHandler.UpdateClass(c)
}

// GetClass get a Class
func (h *Handler) GetClass(c echo.Context) (err error) {
	return h.ClassHandler.GetClass(c)
}

// GetAllClasss get all Classs
func (h *Handler) GetAllClasses(c echo.Context) (err error) {
	return h.ClassHandler.GetAllClasses(c)
}

// DeleteClass delete a Class
func (h *Handler) DeleteClass(c echo.Context) (err error) {
	return h.ClassHandler.DeleteClass(c)
}

// EnrollClass student enroll a Class
func (h *Handler) EnrollClass(c echo.Context) (err error) {
	return h.ClassHandler.EnrollClass(c)
}

// UnnrollClass student unenroll a Class
func (h *Handler) UnnrollClass(c echo.Context) (err error) {
	return h.ClassHandler.UnenrollClass(c)
}

// ListExercises list all exercise in a Class
func (h *Handler) ListExercises(c echo.Context) (err error) {
	return h.ClassHandler.ListExercises(c)
}

func (h *Handler) ListClassUsers(c echo.Context) (err error) {
	return h.ClassHandler.ListClassUsers(c)
}
