package schema

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

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"schema_id": id,
	})
}

func (h *Handler) List(c *gin.Context) {
	schemas, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch schemas",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"schemas": schemas,
	})
}

func (h *Handler) GetLatest(c *gin.Context) {
	idParam := c.Param("schema_id")

	schemaID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid schema id",
		})
		return
	}

	version, snapshot, err := h.service.GetLatest(
		c.Request.Context(),
		schemaID,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "schema not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"version":  version,
		"snapshot": snapshot,
	})
}

func (h *Handler) ListVersions(c *gin.Context) {
	idParam := c.Param("schema_id")

	schemaID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid schema id"})
		return
	}

	versions, err := h.service.ListVersions(c.Request.Context(), schemaID)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to fetch versions"})
		return
	}

	c.JSON(200, gin.H{"versions": versions})
}
