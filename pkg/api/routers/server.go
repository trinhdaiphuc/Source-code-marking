package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/handlers"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/middlewares"
)

func Routing(e *echo.Echo, h *handlers.Handler) {
	landingPage(e, h)
	users(e, h)
	classes(e, h)
	exercises(e, h)
	files(e, h)
	comments(e, h)
	webSocketNotification(e, h)
}

func landingPage(e *echo.Echo, h *handlers.Handler) {
	e.GET("/", h.LandingPage)
	e.GET("/health-check", h.HealthCheck)
	e.GET("/health_check", h.HealthCheck)
}

func users(e *echo.Echo, h *handlers.Handler) {
	e.POST("/api/v1/users/signup", h.Signup)
	e.POST("/api/v1/users/signin", h.Signin)
	e.GET("/api/v1/users/profile", h.Profile)
	e.GET("/api/v1/users/confirmation", h.ValidateUser)
	e.GET("/api/v1/users/password", h.ForgetPassword)
	e.POST("/api/v1/users/password", h.ResetPassword)
	e.PUT("/api/v1/users/password", h.ChangePassword)
	e.GET("/api/v1/users/:id", h.GetUser)
	e.GET("/api/v1/users", h.GetAllUsers, middlewares.IsAdmin)
	e.POST("/api/v1/users", h.CreateUser, middlewares.IsAdmin)
	e.PUT("/api/v1/users/:id", h.UpdateUser)
	e.DELETE("/api/v1/users/:id", h.DeleteUser, middlewares.IsAdmin)
	e.GET("/api/v1/users/:id/classes", h.ListClasses)
	e.GET("/api/v1/users/statistic", h.UserStatistic, middlewares.IsAdmin)
}

func classes(e *echo.Echo, h *handlers.Handler) {
	e.POST("/api/v1/classes", h.CreateClass, middlewares.IsTeacher)
	e.GET("/api/v1/classes/:id", h.GetClass)
	e.GET("/api/v1/classes", h.GetAllClasses)
	e.PUT("/api/v1/classes/:id", h.UpdateClass, middlewares.IsAdminOrTeacher)
	e.DELETE("/api/v1/classes/:id", h.DeleteClass, middlewares.IsAdminOrTeacher)
	e.POST("/api/v1/classes/:id/enroll", h.EnrollClass, middlewares.IsStudent)
	e.PUT("/api/v1/classes/:id/enroll", h.UnnrollClass, middlewares.IsStudent)
	e.GET("/api/v1/classes/:id/exercises", h.ListExercises)
	e.GET("/api/v1/classes/:id/users", h.ListClassUsers, middlewares.IsTeacherOrStudent)
	e.GET("/api/v1/classes/statistic", h.ClassStatistic, middlewares.IsAdmin)
}

func exercises(e *echo.Echo, h *handlers.Handler) {
	e.POST("/api/v1/exercises", h.CreateExercise, middlewares.IsTeacher)
	e.GET("/api/v1/exercises/:id", h.GetExercise)
	e.GET("/api/v1/exercises", h.GetAllExercises)
	e.PUT("/api/v1/exercises/:id", h.UpdateExercise, middlewares.IsAdminOrTeacher)
	e.DELETE("/api/v1/exercises/:id", h.DeleteExercise, middlewares.IsAdminOrTeacher)
	e.GET("/api/v1/exercises/:id/files", h.ListFiles)
	e.GET("/api/v1/exercises/statistic", h.ExerciseStatistic, middlewares.IsAdmin)
}

func files(e *echo.Echo, h *handlers.Handler) {
	e.POST("/api/v1/files", h.CreateFile, middlewares.IsStudent)
	e.GET("/api/v1/files/:id", h.GetFile)
	e.GET("/api/v1/files", h.GetAllFiles, middlewares.IsAdmin)
	e.PUT("/api/v1/files/:id", h.UpdateFile, middlewares.IsAdminOrStudent)
	e.PATCH("/api/v1/files/:id", h.MarkFile, middlewares.IsTeacher)
	e.DELETE("/api/v1/files/:id", h.DeleteFile, middlewares.IsAdminOrStudent)
	e.GET("/api/v1/files/:id/comments", h.ListComments)
	e.GET("/api/v1/files/statistic", h.FileStatistic, middlewares.IsAdmin)
}

func comments(e *echo.Echo, h *handlers.Handler) {
	e.POST("/api/v1/comments", h.CreateComment, middlewares.IsTeacher)
	e.GET("/api/v1/comments/:id", h.GetComment)
	e.GET("/api/v1/comments", h.GetAllComments)
	e.PUT("/api/v1/comments/:id", h.UpdateComment, middlewares.IsAdminOrTeacher)
	e.DELETE("/api/v1/comments/:id", h.DeleteComment, middlewares.IsAdminOrTeacher)
}

func webSocketNotification(e *echo.Echo, h *handlers.Handler) {
	e.GET("/api/v1/ws", h.WebsocketNotification)
	e.GET("/api/v1/notifications", h.GetAllNotifications)
	e.PATCH("/api/v1/notifications/:id", h.MarkReadNotification)
}
