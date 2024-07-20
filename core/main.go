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

	///             PLAYER             \\\

	// GET \\
	r.GET("/players", controllers.GetAllPlayers)
	r.GET("/player/:id", controllers.GetOnePlayer)
	r.GET("/player/:id/stats", controllers.GetPlayerStats)
	r.GET("/player/:id/equipment", controllers.GetPlayerEquipment)
	r.GET("/player/:id/inventory", controllers.GetPlayerInventory)
	r.GET("/player/:id/pets", controllers.GetPlayerPets)
	r.GET("/player/:id/guild", controllers.GetPlayerGuild)
	r.GET("/player/:id/skill", controllers.GetPlayerSkill)
	// POST \\
	r.POST("player/create", controllers.CreatePlayer)
	r.POST("/player/:id/inventory", controllers.AddItemToPlayerInventory)
	r.POST("/player/:id/pets", controllers.AddPetToPlayer)
	r.POST("/player/:id/skill", controllers.AddSkillToPlayer)
	// PATCH \\
	r.PATCH("/player/:id/stats", controllers.UpdatePlayerStats)
	r.PATCH("/player/:id/equipment", controllers.UpdatePlayerEquipment)
	r.PATCH("/player/:id/inventory", controllers.UpdatePlayerInventory)
	r.PATCH("/player/:id/pets", controllers.UpdatePlayerPets)
	r.PATCH("/player/:id/skills", controllers.UpdatePlayerSkills)
	r.PATCH("/player", controllers.UpdatePlayer)
	r.PATCH("/player/location", controllers.UpdatePlayerLocation)
	//DELETE\\
	r.DELETE("/player/:id", controllers.DeletePlayer)
	r.DELETE("/player/:id/inventory", controllers.DeletePlayerItemInInventory)
	r.DELETE("/player/:id/pets", controllers.DeletePlayerPets)
	r.DELETE("/player/:id/skill", controllers.DeletePlayerSkill)

	///             CREATURE             \\\

	// GET \\
	r.GET("/creatures", controllers.GetAllCreatures)
	r.GET("/creature/:id", controllers.GetOneCreature)
	r.GET("/creature/:id/spawn", controllers.GetCreatureSpawn)
	r.GET("/creature/:id/skill", controllers.GetCreatureSkill)
	// POST \\
	r.POST("/creature/create", controllers.CreateCreature)
	r.POST("/creature/spawn", controllers.AddCreatureSpawn)
	r.POST("/creature/skill", controllers.AddCreatureSkill)
	r.POST("/creature/loot", controllers.AddCreatureLoot)
	// PATCH \\
	r.PATCH("/creature/:id", controllers.UpdateCreature)
	r.PATCH("/creature/:id/spawn", controllers.UpdateCreatureSpawn)
	r.PATCH("/creature/:id/skill", controllers.UpdateCreatureSkill)
	r.PATCH("/creature/:id/loot", controllers.UpdateCreatureLoot)
	//DELETE\\
	r.DELETE("/creature/:id", controllers.DeleteCreature)
	r.DELETE("/creature/:id/spawn", controllers.DeleteCreatureSpawn)
	r.DELETE("/creature/:id/skill", controllers.DeleteCreatureSkill)
	r.DELETE("/creature/:id/loot", controllers.DeleteCreatureLoot)

	///             LOCATIONS             \\\
	///             ITEMS             \\\
	///             QUEST             \\\
	///             JOB             \\\
	///             RACE             \\\
	///             RESOURCE             \\\
	///             ACTION             \\\
	///             GUILD             \\\
	///             SKILL             \\\

	err = r.Run()
	if err != nil {
		panic(err)
	}
}
