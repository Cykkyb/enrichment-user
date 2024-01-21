package handler

import (
	"database/sql"
	"enrichment-user-info/internal/entity"
	repository "enrichment-user-info/internal/repository/postgres"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) addUsersHandler(c *gin.Context) {
	h.LogRequest("add user", c)

	var user entity.User

	if err := c.BindJSON(&user); err != nil {
		h.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.services.Enrichment.CreateUser(&user); err != nil {
		h.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, user)
	h.log.Info("add user",
		slog.Any("data", user),
	)
}

func (h *Handler) updateUserHandler(c *gin.Context) {
	h.LogRequest("update user", c)

	var user entity.User

	if err := c.BindJSON(&user); err != nil {
		h.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Id, _ = strconv.Atoi(c.Param("id"))

	if err := h.services.Enrichment.UpdateUser(&user); err != nil {
		h.log.Error(err.Error())
		if err == sql.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": repository.ErrUserNotFound.Error()})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.IndentedJSON(http.StatusOK, user)
	h.log.Info("update user", user)
}

func (h *Handler) deleteUserHandler(c *gin.Context) {
	h.LogRequest("delete user", c)

	id, _ := strconv.Atoi(c.Param("id"))

	if err := h.services.Enrichment.DeleteUser(id); err != nil {
		h.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "user deleted"})
	h.log.Info("delete user",
		slog.Int("id", id),
	)
}

func (h *Handler) getUserHandler(c *gin.Context) {
	h.LogRequest("get user", c)

	var filter entity.UserFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("get user with filter",
		slog.Any("filter", filter),
	)

	users, err := h.services.Enrichment.GetUsersWithFilter(&filter)
	if err != nil {
		h.log.Error(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func (h *Handler) LogRequest(message string, c *gin.Context) {
	h.log.Info("Request:"+message,
		slog.String("ip", c.ClientIP()),
		slog.String("time", time.Now().Format("2006-01-02 15:04:05")),
		slog.String("method", c.Request.Method),
	)
}
