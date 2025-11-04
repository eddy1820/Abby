package api

import (
	_ "Abby/docs" // 這裡會引入自動生成的 swagger 文檔

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Simple Storage API
// @version 1.0
// @description 這是一個簡單的智能合約 API 服務
// @host localhost:8081
// @BasePath /api/v1
// @schemes http
func SetupRouter(handler *StorageHandler) *gin.Engine {
	r := gin.Default()

	// API v1
	v1 := r.Group("/api/v1")
	{
		storage := v1.Group("/storage")
		{
			storage.GET("/value", handler.GetValue)
			storage.POST("/value", handler.SetValue)
		}
	}

	// Swagger 文檔
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
