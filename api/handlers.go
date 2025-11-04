package api

import (
	"math/big"
	"net/http"

	"Abby/contracts"

	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	interactor *contracts.ContractInteractor
}

func NewStorageHandler(interactor *contracts.ContractInteractor) *StorageHandler {
	return &StorageHandler{
		interactor: interactor,
	}
}

// GetValue godoc
// @Summary 獲取存儲的值
// @Description 從智能合約中獲取當前存儲的值
// @Tags storage
// @Accept json
// @Produce json
// @Success 200 {object} object{value=string} "成功返回存儲的值"
// @Failure 500 {object} object{error=string} "內部錯誤"
// @Router /storage/value [get]
func (h *StorageHandler) GetValue(c *gin.Context) {
	value, err := h.interactor.GetValue()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"value": value.String(),
	})
}

// SetValueRequest 設置值的請求結構
type SetValueRequest struct {
	Value string `json:"value" example:"42" binding:"required"`
}

// SetValue godoc
// @Summary 設置新的值
// @Description 在智能合約中設置新的值
// @Tags storage
// @Accept json
// @Produce json
// @Param request body SetValueRequest true "要設置的新值"
// @Success 200 {object} object{message=string} "成功設置值"
// @Failure 400 {object} object{error=string} "請求格式錯誤"
// @Failure 500 {object} object{error=string} "內部錯誤"
// @Router /storage/value [post]
func (h *StorageHandler) SetValue(c *gin.Context) {
	var request struct {
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// 將字符串轉換為 big.Int
	value := new(big.Int)
	value, ok := value.SetString(request.Value, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid number format",
		})
		return
	}

	err := h.interactor.SetValue(value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Value updated successfully",
	})
}
