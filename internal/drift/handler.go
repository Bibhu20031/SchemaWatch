package drift

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c *gin.Context) {
	idParam := c.Param("schema_id")

	schemaID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schema id"})
		return
	}

	events, err := h.service.ListBySchema(
		c.Request.Context(),
		schemaID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch drift events",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"drifts": events,
	})
}
