package handler

import (
	"enrichment-user-info/internal/service"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type Handler struct {
	services *service.Service
	log      *slog.Logger
}

func NewHandler(services *service.Service, log *slog.Logger) *Handler {
	return &Handler{
		services: services,
		log:      log,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			users := v1.Group("/users")
			{
				users.POST("/", h.addUsersHandler)
				users.GET("/", h.getUserHandler)
				users.PUT("/:id", h.updateUserHandler)
				users.DELETE("/:id", h.deleteUserHandler)
			}
		}
	}

	return router
}
