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

// DeviceHandler HTTPリクエストを処理し、DeviceUsecaseを呼び出す.
type DeviceHandler struct {
	uc usecase.DeviceUsecase
}

// NewDeviceHandler DeviceHandlerの新しいインスタンスを生成.
func NewDeviceHandler(uc usecase.DeviceUsecase) *DeviceHandler {
	return &DeviceHandler{uc: uc}
}

// CreateDevice POST /devices - 新しいデバイスを作成.
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

// GetDevice GET /devices/:id - 特定のデバイスを取得.
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

// ListDevices GET /devices - 全てのデバイスを取得.
func (h *DeviceHandler) ListDevices(c *gin.Context) {
	outputs, err := h.uc.ListDevices(c.Request.Context())
	if err != nil {
		log.Printf("failed to list devices: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})

		return
	}

	// 空のスライスの場合でも nil ではなく [] を返すようにする
	if outputs == nil {
		outputs = []*usecase.DeviceOutput{}
	}

	c.JSON(http.StatusOK, outputs)
}

// UpdateDevice PUT /devices/:id - 特定のデバイスを更新.
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

	input.ID = id // URLから取得したIDをインプットに設定

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

// DeleteDevice DELETE /devices/:id - 特定のデバイスを削除.
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
