package routes

import (
	"Slackluky/workfly/internal/handler"
	"Slackluky/workfly/internal/middleware"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
			auth.POST("/forgot-password", handler.RequestPasswordReset)
			auth.POST("/reset-password", handler.ConfirmPasswordReset)
		}

		users := api.Group("/users")
		users.Use(middleware.JWTAuthMiddleware)
		{
			users.GET("/me", handler.GetMe)
			users.PATCH("/me", handler.UpdateMe)
			users.DELETE("/me", handler.DeleteMe)

			users.GET("", handler.ListUsers)
			users.GET("/:id", handler.GetUser)
			users.PATCH("/:id", handler.UpdateUser)
			users.DELETE("/:id", handler.DeleteUser)
		}

	}
}
