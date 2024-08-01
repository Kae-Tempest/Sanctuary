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

	// TODO : Uniformiser les routes

	///             USER             \\\
	users := r.Group("/users")
	// GET \\
	users.GET("/", controllers.GetUsers)
	users.GET("/:id", controllers.GetUserByID)
	users.GET("/e/:email", controllers.GetUserByEmail)
	// POST \\
	users.POST("/register", controllers.Register)
	users.POST("/login", controllers.Login)

	///             PLAYER             \\\ TODO : Player -> Character
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
	//characters.PATCH("/player/:id/pets", controllers.UpdateCharactersPets) // TODO: Review Update
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

	///             LOCATIONS             \\\ TODO : Definir des actions possible / zones
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
	// PATCH \\
	locations.PATCH("/:id", controllers.UpdateLocation)
	//DELETE\\
	locations.DELETE("/:id", controllers.DeleteLocation)

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
	///             SKILL             \\\ // TODO : Gestion d'assignation de skill
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
	races.GET("/races")
	races.GET("/race/:id")
	// POST \\
	races.POST("/race/create")
	// PATCH \\
	races.PATCH("/race/:id")
	//DELETE\\
	races.DELETE("/race/:id")
	///             JOB             \\\
	jobs := r.Group("jobs")
	// GET \\
	jobs.GET("/jobs")
	jobs.GET("/job/:id")
	jobs.GET("/job/:id/skills")
	jobs.GET("/job/:id/skill/:skill")
	// POST \\
	jobs.POST("/job/create")
	jobs.POST("/job/skill/create")
	// PATCH \\
	jobs.PATCH("/job/:id")
	jobs.PATCH("/job/skill/:id")
	//DELETE\\
	jobs.DELETE("/job/:id")
	jobs.DELETE("/job/skill/:id")

	///             ACTION LOGS            \\\
	logs := r.Group("/logs")
	// GET \\
	logs.GET("/actions")
	logs.GET("/actions/:player")
	logs.GET("/actions/:type")
	// POST \\
	logs.POST("/action/create")
	// PATCH \\
	logs.PATCH("/action/:id")
	//DELETE\\
	logs.DELETE("/action/:id")
	///             GUILD             \\\
	guilds := r.Group("/guilds")
	// GET \\
	guilds.GET("/guilds")
	guilds.GET("/guild/:id")
	guilds.GET("/guild/owner/:owner")
	guilds.GET("guild/:id/members")
	// POST \\
	guilds.POST("/guild/create")
	guilds.POST("/guild/:id/invite/:player")
	// PATCH \\
	guilds.PATCH("/guild/:id")
	//DELETE\\
	guilds.DELETE("/guild/:id")
	guilds.DELETE("/guild/:id/eject/:player")
	///             QUEST             \\\
	quests := r.Group("/quests")
	// GET \\
	quests.GET("/quests")
	quests.GET("/quest/:id")
	quests.GET("/quest/type/:type")
	quests.GET("/quest/rank/:rank")
	// POST \\
	quests.POST("/quest/create")
	// PATCH \\
	quests.PATCH("/quest/:id")
	//DELETE\\
	quests.DELETE("/quest/:id")

	///             RESOURCE             \\\ TODO : Revoir les resources ( + DB )
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

	err = r.Run(":8000")
	if err != nil {
		panic(err)
	}
}
