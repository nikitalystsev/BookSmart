package handlers

import "github.com/gin-gonic/gin"

func (h *Handler) initAdminRoutes(api *gin.RouterGroup) {
	students := api.Group("/admins")
	{
		students.POST("/sign-up")
		students.POST("/sign-in")
		students.POST("/auth/refresh")
	}
}
