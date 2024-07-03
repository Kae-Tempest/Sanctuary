package main

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"sanctuary-api/controllers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
		c.Done()
	})

	r.GET("/players", controllers.GetAllPlayers)
	r.GET("/player/:id", controllers.GetOnePlayer)
	r.GET("/player/:id/stats", controllers.GetPlayerStats)
	r.GET("/player/:id/equipment", controllers.GetPlayerEquipment)
	r.GET("/player/:id/inventory", controllers.GetPlayerInventory)
	r.GET("/player/:id/pets", controllers.GetPlayerPets)
	r.GET("/player/:id/guild", controllers.GetPlayerGuild)
	r.GET("/player/:id/skill", controllers.GetPlayerSkill)
	r.POST("player/:id", controllers.CreatePlayer)
	r.POST("/player/:id/inventory", controllers.AddItemToPlayerInventory)
	r.POST("/player/:id/pets", controllers.AddPetToPlayer)
	r.POST("/player/:id/skill", controllers.AddSkillToPlayer)

	err = r.Run()
	if err != nil {
		panic(err)
	}
}
