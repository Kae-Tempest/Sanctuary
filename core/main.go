package main

import (
	"log/slog"
	"net/http"

	"sanctuary-api/controllers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	})

	///             USER             \\\
	// GET \\
	r.GET("/users", controllers.GetUsers)
	r.GET("/user/:id", controllers.GetUserByID)
	r.GET("/user/e/:email", controllers.GetUserByEmail)
	// POST \\
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
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
	r.POST("/player/create", controllers.CreatePlayer)
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
	r.GET("/mobs", controllers.GetAllMobs)
	r.GET("/mob/:id", controllers.GetOneMob)
	r.GET("/mob/:id/spawn", controllers.GetMobSpawn)
	r.GET("/mob/:id/skill", controllers.GetMobSkill)
	// POST \\
	r.POST("/mob/create", controllers.CreateMob)
	r.POST("/mob/spawn", controllers.AddMobSpawn)
	r.POST("/mob/skill", controllers.AddMobSkill)
	r.POST("/mob/loot", controllers.AddMobLoot)
	// PATCH \\
	r.PATCH("/mob/:id", controllers.UpdateMob)
	r.PATCH("/mob/:id/spawn", controllers.UpdateMobSpawn)
	r.PATCH("/mob/:id/skill", controllers.UpdateMobSkill)
	r.PATCH("/mob/:id/loot", controllers.UpdateMobLoot)
	//DELETE\\
	r.DELETE("/mob/:id", controllers.DeleteMob)
	r.DELETE("/mob/:id/spawn", controllers.DeleteMobSpawn)
	r.DELETE("/mob/:id/skill", controllers.DeleteMobSkill)
	r.DELETE("/mob/:id/loot", controllers.DeleteMobLoot)

	///             LOCATIONS             \\\
	// GET \\
	r.GET("/locations", controllers.GetLocations)
	r.GET("/location/:id", controllers.GetLocationByID)
	r.GET("/location/:id/players", controllers.GetPlayersByLocation)
	r.GET("/location/:id/mobs", controllers.GetCreaturesByLocation)
	r.GET("/location/:id/resources", controllers.GetResourcesByLocation)
	r.GET("/location/:id/loots")
	// POST \\
	r.POST("/location", controllers.CreateLocation)
	// PATCH \\
	r.PATCH("/location/:id", controllers.UpdateLocation)
	//DELETE\\
	r.DELETE("/location/:id", controllers.DeleteLocation)

	///             ITEMS             \\\
	// GET \\
	r.GET("/items", controllers.GetItems)
	r.GET("/item/:id", controllers.GetItemByID)
	r.GET("/item/type/:type", controllers.GetItemByType)
	// POST \\
	r.POST("/item/create", controllers.CreateItem)
	// PATCH \\
	r.PATCH("/item/:id", controllers.UpdateItem)
	r.PATCH("/item/:id/stat", controllers.UpdateItemStat)
	r.PATCH("/item/:id/emplacement", controllers.UpdateItemEmplacement)
	//DELETE\\
	r.DELETE("/item/:id", controllers.DeleteItem)
	///             SKILL             \\\
	// GET \\
	r.GET("/skills", controllers.GetSkills)
	r.GET("/skill/:id", controllers.GetSkillByID)
	r.GET("/skill/type/:type", controllers.GetSkillByType)
	// POST \\
	r.POST("/skill/create", controllers.CreateSkill)
	// PATCH \\
	r.PATCH("/skill/:id/info", controllers.UpdateSkillInfo)
	r.PATCH("/skill/:id/stats", controllers.UpdateSkillStats)
	//DELETE\\
	r.DELETE("/skill/:id", controllers.DeleteSkill)
	///             RACE             \\\
	// GET \\
	r.GET("/races")
	r.GET("/race/:id")
	// POST \\
	r.POST("/race/create")
	// PATCH \\
	r.PATCH("/race/:id")
	//DELETE\\
	r.DELETE("/race/:id")
	///             JOB             \\\
	// GET \\
	r.GET("/jobs")
	r.GET("/job/:id")
	r.GET("/job/:id/skills")
	r.GET("/job/:id/skill/:skill")
	// POST \\
	r.POST("/job/create")
	r.POST("/job/skill/create")
	// PATCH \\
	r.PATCH("/job/:id")
	r.PATCH("/job/skill/:id")
	//DELETE\\
	r.DELETE("/job/:id")
	r.DELETE("/job/skill/:id")
	///             RESOURCE             \\\
	// GET \\
	r.GET("/resources")
	r.GET("/resource/:id")
	r.GET("/resource/type/:type")
	r.GET("/resource/location/:id")
	// POST \\
	r.POST("/resource/create")
	r.POST("/resource/type/create")
	// PATCH \\
	r.PATCH("/resource/:id")
	r.PATCH("/resource/type/:id")
	//DELETE\\
	r.DELETE("/resource/:id")
	r.DELETE("/resource/type/:id")
	///             ACTION             \\\
	// GET \\
	r.GET("/actions")
	r.GET("/actions/:player")
	r.GET("/actions/:type")
	// POST \\
	r.POST("/action/create")
	// PATCH \\
	r.PATCH("/action/:id")
	//DELETE\\
	r.DELETE("/action/:id")
	///             GUILD             \\\
	// GET \\
	r.GET("/guilds")
	r.GET("/guild/:id")
	r.GET("/guild/owner/:owner")
	r.GET("guild/:id/members")
	// POST \\
	r.POST("/guild/create")
	r.POST("/guild/:id/invite/:player")
	// PATCH \\
	r.PATCH("/guild/:id")
	//DELETE\\
	r.DELETE("/guild/:id")
	r.DELETE("/guild/:id/eject/:player")
	///             QUEST             \\\
	// GET \\
	r.GET("/quests")
	r.GET("/quest/:id")
	r.GET("/quest/type/:type")
	r.GET("/quest/rank/:rank")
	// POST \\
	r.POST("/quest/create")
	// PATCH \\
	r.PATCH("/quest/:id")
	//DELETE\\
	r.DELETE("/quest/:id")

	err = r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
