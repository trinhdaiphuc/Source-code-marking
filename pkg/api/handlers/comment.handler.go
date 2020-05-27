package handlers

import "github.com/labstack/echo/v4"

/////////////////////////////////////////////////////////
// TODO Comment handler
/////////////////////////////////////////////////////////

// CreateComment create a comment in a file
func (h *Handler) CreateComment(c echo.Context) (err error) {
	return h.CommentHandler.CreateComment(c)
}

// UpdateComment update a comment in a file
func (h *Handler) UpdateComment(c echo.Context) (err error) {
	return h.CommentHandler.UpdateComment(c)
}

// GetComment get a comment in a file
func (h *Handler) GetComment(c echo.Context) (err error) {
	return h.CommentHandler.GetComment(c)
}

// GetAllComments get all comments in a file
func (h *Handler) GetAllComments(c echo.Context) (err error) {
	return h.CommentHandler.GetAllComments(c)
}

// DeleteComment delete a comment in a file
func (h *Handler) DeleteComment(c echo.Context) (err error) {
	return h.CommentHandler.DeleteComment(c)
}
