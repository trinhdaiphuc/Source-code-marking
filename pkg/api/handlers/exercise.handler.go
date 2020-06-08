package handlers

import "github.com/labstack/echo/v4"

/////////////////////////////////////////////////////////
// TODO Exercise handler
/////////////////////////////////////////////////////////

// CreateExercise create a Exercise
func (h *Handler) CreateExercise(c echo.Context) (err error) {
	return h.ExerciseHandler.CreateExercise(c)
}

// UpdateExercise update a Exercise
func (h *Handler) UpdateExercise(c echo.Context) (err error) {
	return h.ExerciseHandler.UpdateExercise(c)
}

// GetExercise get a Exercise
func (h *Handler) GetExercise(c echo.Context) (err error) {
	return h.ExerciseHandler.GetExercise(c)
}

// GetAllExercises get all Exercises
func (h *Handler) GetAllExercises(c echo.Context) (err error) {
	return h.ExerciseHandler.GetAllExercises(c)
}

// DeleteExercise delete a Exercise
func (h *Handler) DeleteExercise(c echo.Context) (err error) {
	return h.ExerciseHandler.DeleteExercise(c)
}

// ListFiles list all file in a Exercise
func (h *Handler) ListFiles(c echo.Context) (err error) {
	return h.ExerciseHandler.ListFiles(c)
}

func (h *Handler) ExerciseStatistic(c echo.Context) (err error) {
	return h.ExerciseHandler.ExerciseStatistic(c)
}
