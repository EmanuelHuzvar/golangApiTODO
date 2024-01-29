package config

import (
	"SamkoOfMraz/api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {

	router.Use(handlers.CorsMiddleware())

	router.POST("/login", handlers.LoginHandler)
	router.POST("/logout", handlers.LogoutHandler)
	router.PUT("/task", handlers.EditTaskHandler)
	router.POST("/task", handlers.AddTaskHandler)
	router.DELETE("/task", handlers.RemoveTaskByIDHandler)
	router.GET("/task", handlers.GetTasksHandler)
	router.PUT("/make_user", handlers.MakeUserHandler)
	router.GET("/state", handlers.GetUserStatesHandler)
	router.POST("/state", handlers.UpdateUserStatesHandler)

	router.POST("/estimate", handlers.UpdateTaskEstimatesHandler)
	router.POST("/task_state", handlers.UpdateTaskStateHandler)
	router.POST("/task_position", handlers.UpdateTaskPositionHandler)

}
