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
	users := r.Group("/users")
	// GET \\
	users.GET("/", controllers.GetUsers)
	users.GET("/:id", controllers.GetUserByID)
	users.GET("/e/:email", controllers.GetUserByEmail)
	// POST \\
	users.POST("/register", controllers.Register)
	users.POST("/login", controllers.Login)

	///             PLAYER             \\\
	characters := r.Group("/characters")
	// GET \\
	characters.GET("/", controllers.GetAllCharacters)
	characters.GET("/:id", controllers.GetOneCharacters)
	characters.GET("/:id/stats", controllers.GetCharactersStats)
	characters.GET("/:id/equipment", controllers.GetCharactersEquipment)
	characters.GET("/:id/inventory", controllers.GetCharactersInventory)
	characters.GET("/:id/pets", controllers.GetCharactersPets)
	characters.GET("/:id/guild", controllers.GetCharactersGuild)
	characters.GET("/:id/skill", controllers.GetCharactersSkill)
	// POST \\
	characters.POST("/", controllers.CreateCharacters)
	characters.POST("/:id/inventory", controllers.AddItemToCharactersInventory)
	characters.POST("/:id/pets", controllers.AddPetToCharacters)
	characters.POST("/:id/skill", controllers.AddSkillToCharacters)
	// PATCH \\
	characters.PATCH("/:id/stats", controllers.UpdateCharactersStats)
	characters.PATCH("/:id/equipment", controllers.UpdateCharactersEquipment)
	characters.PATCH("/:id/inventory", controllers.UpdateCharactersInventory)
	characters.PATCH("/:id/skills", controllers.UpdateCharactersSkills)
	characters.PATCH("/:id", controllers.UpdateCharacters)
	characters.PATCH("/:id/location", controllers.UpdateCharactersLocation)
	//DELETE\\
	characters.DELETE("/:id", controllers.DeleteCharacters)
	characters.DELETE("/:id/inventory", controllers.DeleteCharactersItemInInventory)
	characters.DELETE("/:id/pets", controllers.DeleteCharactersPets)
	characters.DELETE("/:id/skill", controllers.DeleteCharactersSkill)

	///             MOB             \\\
	mobs := r.Group("/mobs")
	// GET \\
	mobs.GET("/", controllers.GetAllMobs)
	mobs.GET("/:id", controllers.GetMobByID)
	mobs.GET("/:id/spawn", controllers.GetMobSpawn)
	mobs.GET("/:id/skill", controllers.GetMobSkill)
	// POST \\
	mobs.POST("/", controllers.CreateMob)
	mobs.POST(":id/spawn", controllers.AddMobSpawn)
	mobs.POST(":id/skill", controllers.AddMobSkill)
	mobs.POST(":id/loot", controllers.AddMobLoot)
	// PATCH \\
	mobs.PATCH("/:id", controllers.UpdateMob)
	mobs.PATCH("/:id/spawn", controllers.UpdateMobSpawn)
	mobs.PATCH("/:id/skill", controllers.UpdateMobSkill)
	mobs.PATCH("/:id/loot", controllers.UpdateMobLoot)
	//DELETE\\
	mobs.DELETE("/:id", controllers.DeleteMob)
	mobs.DELETE("/:id/spawn", controllers.DeleteMobSpawn)
	mobs.DELETE("/:id/skill", controllers.DeleteMobSkill)
	mobs.DELETE("/:id/loot", controllers.DeleteMobLoot)

	///             LOCATIONS             \\\
	locations := r.Group("/")
	// GET \\
	locations.GET("/", controllers.GetLocations)
	locations.GET("/:id", controllers.GetLocationByID)
	locations.GET("/:id/players", controllers.GetCharactersByLocation)
	locations.GET("/:id/mobs", controllers.GetCreaturesByLocation)
	locations.GET("/:id/resources", controllers.GetResourcesByLocation)
	locations.GET("/:id/actions")
	locations.GET("/:id/loots")
	// POST \\
	locations.POST("/", controllers.CreateLocation)
	locations.POST("/:id/actions")
	// PATCH \\
	locations.PATCH("/:id", controllers.UpdateLocation)
	//DELETE\\
	locations.DELETE("/:id", controllers.DeleteLocation)
	locations.DELETE("/:id/actions")

	///             ITEMS             \\\
	items := r.Group("/items")
	// GET \\
	items.GET("/", controllers.GetItems)
	items.GET("/:id", controllers.GetItemByID)
	items.GET("/type/:type", controllers.GetItemByType)
	// POST \\
	items.POST("/", controllers.CreateItem)
	// PATCH \\
	items.PATCH("/:id", controllers.UpdateItem)
	items.PATCH("/:id/stat", controllers.UpdateItemStat)
	items.PATCH("/:id/emplacement", controllers.UpdateItemEmplacement)
	//DELETE\\
	items.DELETE("/:id", controllers.DeleteItem)
	///             SKILL             \\\
	skills := r.Group("/skills")
	// GET \\
	skills.GET("/", controllers.GetSkills)
	skills.GET("/:id", controllers.GetSkillByID)
	skills.GET("/type/:type", controllers.GetSkillByType)
	// POST \\
	skills.POST("/", controllers.CreateSkill)
	// PATCH \\
	skills.PATCH("/:id/info", controllers.UpdateSkillInfo)
	skills.PATCH("/:id/stats", controllers.UpdateSkillStats)
	//DELETE\\
	skills.DELETE("/:id", controllers.DeleteSkill)
	///             RACE             \\\
	races := r.Group("/races")
	// GET \\
	races.GET("/")
	races.GET("//:id")
	// POST \\
	races.POST("/")
	// PATCH \\
	races.PATCH("/:id")
	//DELETE\\
	races.DELETE("/:id")
	///             JOB             \\\
	jobs := r.Group("jobs")
	// GET \\
	jobs.GET("")
	jobs.GET("/:id")
	jobs.GET("/:id/skills")
	jobs.GET("/:id/skill/:skill")
	jobs.GET("/:id/skill_tree")
	// POST \\
	jobs.POST("/")
	jobs.POST("/skill")
	// PATCH \\
	jobs.PATCH("/:id")
	jobs.PATCH("/skill/:id")
	//DELETE\\
	jobs.DELETE("/:id")
	jobs.DELETE("/skill/:id")

	///             ACTION LOGS            \\\
	logs := r.Group("/logs")
	// GET \\
	logs.GET("/")
	logs.GET("/:player")
	logs.GET("/:type")
	// POST \\
	logs.POST("/")
	// PATCH \\
	logs.PATCH("/:id")
	//DELETE\\
	logs.DELETE("/:id")
	///             GUILD             \\\
	guilds := r.Group("/guilds")
	// GET \\
	guilds.GET("/")
	guilds.GET("/:id")
	guilds.GET("/owner/:owner")
	guilds.GET("/:id/members")
	// POST \\
	guilds.POST("/")
	guilds.POST("/:id/invite/:player")
	// PATCH \\
	guilds.PATCH("/:id")
	//DELETE\\
	guilds.DELETE("/:id")
	guilds.DELETE("/guild/:id/eject/:player")
	///             QUEST             \\\
	quests := r.Group("/quests")
	// GET \\
	quests.GET("/")
	quests.GET("/:id")
	quests.GET("/type/:type")
	quests.GET("/rank/:rank")
	// POST \\
	quests.POST("/")
	// PATCH \\
	quests.PATCH("/:id")
	//DELETE\\
	quests.DELETE("/:id")

	///             RESOURCE             \\\
	resources := r.Group("/resources")
	// GET \\
	resources.GET("/")
	resources.GET("/:id")
	resources.GET("/type/:type")
	resources.GET("/:id/location")
	// POST \\
	resources.POST("/create")
	// PATCH \\
	resources.PATCH("/:id")
	//DELETE\\
	resources.DELETE("/:id")

	err = r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
