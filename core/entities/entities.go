package entities

import "time"

type Player struct {
	ID            int    `json:"id"`
	Email         string `json:"email"`
	Username      string `json:"username"`
	RaceID        int    `json:"race_id"`
	JobID         int    `json:"job_id"`
	Exp           int    `json:"exp"`
	Po            int    `json:"po"`
	Level         int    `json:"level"`
	GuildID       int    `json:"guild_id"` // 0 = no guild
	InventorySize int    `json:"inventorySize"`
	LocationId    int    `json:"locationId"`
}

type Inventory struct {
	PlayerID int `json:"player_id"`
	ItemID   int `json:"item_id"`
	Quantity int `json:"quantity"`
}

type Job struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"` // description of the job
	Strength     int    `json:"strength"`
	Constitution int    `json:"constitution"`
	Dexterity    int    `json:"dexterity"`
	Intelligence int    `json:"intelligence"`
	Mana         int    `json:"mana"`
	Stamina      int    `json:"stamina"`
	Wisdom       int    `json:"wisdom"`
	Charisma     int    `json:"charisma"`
}

type Race struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"` // description of the race
	Mana        int    `json:"mana"`
	Stamina     int    `json:"stamina"`
	Wisdom      int    `json:"wisdom"`
	Charisma    int    `json:"charisma"`
}

type Items struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Type             int    `json:"type"` // 0 = Equipable, 1 = Consomable, 2 = Quest, 3 = Resources
	Strength         int    `json:"strength"`
	Constitution     int    `json:"constitution"`
	Mana             int    `json:"mana"`
	Stamina          int    `json:"stamina"`
	Dexterity        int    `json:"dexterity"`
	Intelligence     int    `json:"intelligence"`
	Wisdom           int    `json:"wisdom"`
	Charisma         int    `json:"charisma"`
	EnchantmentLevel int    `json:"enchantmentLevel"`
}

type Guild struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Members []int  `json:"members"`
	Owner   string `json:"owner"`
}

type GuildMembers struct {
	ID       int `json:"id"`
	GuildId  int `json:"guild_id"`
	PlayerID int `json:"player_id"`
}

type Skill struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"` // physic or magic
	Description  string `json:"description"`
	Strength     int    `json:"strength"`
	Constitution int    `json:"constitution"`
	Mana         int    `json:"mana"`
	Stamina      int    `json:"stamina"`
	Dexterity    int    `json:"dexterity"`
	Intelligence int    `json:"intelligence"`
	Wisdom       int    `json:"wisdom"`
	Charisma     int    `json:"charisma"`
}

type JobSkill struct {
	ID           int    `json:"id"`
	JobID        int    `json:"jobID"`
	Name         string `json:"name"`
	Type         string `json:"type"` // physic or magic
	Description  string `json:"description"`
	Strength     int    `json:"strength"`
	Constitution int    `json:"constitution"`
	Mana         int    `json:"mana"`
	Stamina      int    `json:"stamina"`
	Dexterity    int    `json:"dexterity"`
	Intelligence int    `json:"intelligence"`
	Wisdom       int    `json:"wisdom"`
	Charisma     int    `json:"charisma"`
}

type UserSkill struct {
	PlayerId int
	SkillID  int
}

type UserJobSkill struct {
	PlayerId int
	SkillID  int
}

type CreatureSkill struct {
	CreatureId int
	SkillID    int
}

type UserPet struct {
	PetID    int `json:"pet_id"`
	PlayerId int `json:"user_id"`
}
type SummonBeast struct {
	ID           int    `json:"id"`
	PlayerId     int    `json:"user_id"`
	Name         string `json:"name"`
	Strength     int    `json:"strength"`
	Constitution int    `json:"constitution"`
	Mana         int    `json:"mana"`
	Stamina      int    `json:"stamina"`
	Dexterity    int    `json:"dexterity"`
	Intelligence int    `json:"intelligence"`
	Wisdom       int    `json:"wisdom"`
	Charisma     int    `json:"charisma"`
}
type Stats struct {
	PlayerId     int `json:"user_id"`
	HP           int `json:"HP"`
	Strength     int `json:"strength"`
	Constitution int `json:"constitution"`
	Mana         int `json:"mana"`
	Stamina      int `json:"stamina"`
	Dexterity    int `json:"dexterity"`
	Intelligence int `json:"intelligence"`
	Charisma     int `json:"charisma"`
	Wisdom       int `json:"wisdom"`
}
type PetsMounts struct {
	ID          int  `json:"id"`
	CreatureID  int  `json:"creature_id"`
	IsMountable bool `json:"is_mountable"`
	Speed       int  `json:"speed"` // 0 = slow, 1 = normal, 2 = fast
}
type Equipment struct {
	PlayerId   int `json:"player_id"`
	Helmet     int `json:"helmet"`
	Chestplate int `json:"chestplate"`
	Leggings   int `json:"leggings"`
	Boots      int `json:"boots"`
	Mainhand   int `json:"mainhand"`
	Offhand    int `json:"offhand"`
	Accesory0  int `json:"accesory_0"`
	Accesory1  int `json:"accesory_1"`
	Accesory2  int `json:"accesory_2"`
	Accesory3  int `json:"accesory_3"`
}

type Creatures struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	IsPet        bool   `json:"is_pet"`
	Level        int    `json:"level"`
	HP           int    `json:"HP"`
	Strength     int    `json:"strength"`
	Constitution int    `json:"constitution"`
	Mana         int    `json:"mana"`
	Stamina      int    `json:"stamina"`
	Dexterity    int    `json:"dexterity"`
	Intelligence int    `json:"intelligence"`
	Charisma     int    `json:"charisma"`
	Wisdom       int    `json:"wisdom"`
}

type CreatureSpawns struct {
	CreatureID    int `json:"creatureID"`
	EmplacementID int `json:"emplacementID"`
	LevelRequired int `json:"levelRequired"`
	SpawnRate     int `json:"spawnRate"`
}

type Quests struct {
	ID          int          `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	IsGroup     bool         `json:"is_group"`
	Difficulty  int          `json:"difficulty"`
	Data        []Objectives `json:"data"`
	Reward      Rewards      `json:"reward"`
}

type Rewards struct {
	Exp  int   `json:"exp"`
	Po   int   `json:"po"`
	Item []int `json:"item"`
}

type Objectives struct {
	Title        string `json:"title"`         // {"objectif": "tuer 10 monstres"}
	Objective    int    `json:"objective"`     // {"track": 0}
	MaxObjective int    `json:"max_objective"` // {"max_track": 10}
}

type Locations struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Type       int    `json:"type"` // 0 = Resource, 1 = Mobs, 2 = City
	IsSafety   bool   `json:"isSafety"`
	Difficulty int    `json:"difficulty"`
	Size       int    `json:"size"`
}

type Resources struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	EmplacementsId   int    `json:"emplacementsId"`
	ResourcesTypeId  int    `json:"resourcesTypeId"`
	QuantitiesPerMin int    `json:"quantitiesPerMin"`
}

type ResourcesType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type PlayerAction struct {
	PlayerId  int       `json:"user_id"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
	EndAt     time.Time `json:"end_at"`
}

type FightOrder struct {
	Name string
	ID   int
}

type HuntAction struct {
	PlayerID  int
	BtnID     string
	MessageID string
	ChannelID string
}
