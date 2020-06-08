package handlers

import "github.com/labstack/echo/v4"

/////////////////////////////////////////////////////////
// TODO File handler
/////////////////////////////////////////////////////////

// CreateFile create a file in an exercise
func (h *Handler) CreateFile(c echo.Context) (err error) {
	return h.FileHandler.CreateFile(c)
}

// UpdateFile update a file in an exercise
func (h *Handler) UpdateFile(c echo.Context) (err error) {
	return h.FileHandler.UpdateFile(c)
}

// MarkFile mark a file in an exercise
func (h *Handler) MarkFile(c echo.Context) (err error) {
	return h.FileHandler.MarkFile(c)
}

// GetFile get a file in an exercise
func (h *Handler) GetFile(c echo.Context) (err error) {
	return h.FileHandler.GetFile(c)
}

// GetAllFiles get all files in an exercise
func (h *Handler) GetAllFiles(c echo.Context) (err error) {
	return h.FileHandler.GetAllFiles(c)
}

// DeleteFile delete a file in an exercise
func (h *Handler) DeleteFile(c echo.Context) (err error) {
	return h.FileHandler.DeleteFile(c)
}

// ListComments list all comments in a file
func (h *Handler) ListComments(c echo.Context) (err error) {
	return h.FileHandler.ListComments(c)
}

func (h *Handler) FileStatistic(c echo.Context) (err error) {
	return h.FileHandler.FileStatistic(c)
}
