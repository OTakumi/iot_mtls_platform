// Package handler provides HTTP handlers for the application.
package handler

import (
	"errors"
	"log"
	"net/http"

	"backend/internal/domain/entity"
	"backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// DeviceHandler handles HTTP requests and calls the DeviceUsecase.
type DeviceHandler struct {
	uc usecase.DeviceUsecase
}

// NewDeviceHandler creates a new instance of DeviceHandler.
func NewDeviceHandler(uc usecase.DeviceUsecase) *DeviceHandler {
	return &DeviceHandler{uc: uc}
}

// CreateDevice handles POST /devices to create a new device.
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var input usecase.CreateDeviceInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})

		return
	}

	output, err := h.uc.CreateDevice(c.Request.Context(), input)
	if err != nil {
		log.Printf("failed to create device: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})

		return
	}

	c.JSON(http.StatusCreated, output)
}

// GetDevice handles GET /devices/:id to retrieve a specific device.
func (h *DeviceHandler) GetDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})

		return
	}

	output, err := h.uc.GetDevice(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrDBFindByID) || errors.Is(err, entity.ErrDeviceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": entity.ErrDeviceNotFound.Error()})

			return
		}

		log.Printf("failed to get device: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})

		return
	}

	c.JSON(http.StatusOK, output)
}

// ListDevices handles GET /devices to retrieve all devices.
func (h *DeviceHandler) ListDevices(c *gin.Context) {
	outputs, err := h.uc.ListDevices(c.Request.Context())
	if err != nil {
		log.Printf("failed to list devices: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})

		return
	}

	// Ensure the response is an empty array `[]` instead of `null` if the slice is empty.
	if outputs == nil {
		outputs = []*usecase.DeviceOutput{}
	}

	c.JSON(http.StatusOK, outputs)
}

// UpdateDevice handles PUT /devices/:id to update a specific device.
func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})

		return
	}

	var input usecase.UpdateDeviceInput

	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})

		return
	}

	input.ID = id // Set the ID from the URL into the input struct.

	output, err := h.uc.UpdateDevice(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, usecase.ErrDBFindByID) || errors.Is(err, entity.ErrDeviceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": entity.ErrDeviceNotFound.Error()})

			return
		}

		log.Printf("failed to update device: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})

		return
	}

	c.JSON(http.StatusOK, output)
}

// DeleteDevice handles DELETE /devices/:id to delete a specific device.
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})

		return
	}

	err = h.uc.DeleteDevice(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, usecase.ErrDBDelete) || errors.Is(err, entity.ErrDeviceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": entity.ErrDeviceNotFound.Error()})

			return
		}

		log.Printf("failed to delete device: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})

		return
	}

	c.Status(http.StatusNoContent)
}
