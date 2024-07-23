package entities

import "time"

type Equipment struct {
	PlayerId       int `json:"player_id"`
	Helmet         int `json:"helmet"`
	Chestplate     int `json:"chestplate"`
	Leggings       int `json:"leggings"`
	Boots          int `json:"boots"`
	Mainhand       int `json:"mainhand"`
	Offhand        int `json:"offhand"`
	AccessorySlot0 int `json:"accessory_slot_0" db:"accessory_slot_0"`
	AccessorySlot1 int `json:"accessory_slot_1" db:"accessory_slot_1"`
	AccessorySlot2 int `json:"accessory_slot_2" db:"accessory_slot_2"`
	AccessorySlot3 int `json:"accessory_slot_3" db:"accessory_slot_3"`
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

type HuntAction struct {
	PlayerID   int       `json:"playerID"`
	LocationID int       `json:"locationID"`
	MobID      int       `json:"mobID"`
	StartAt    time.Time `json:"startAt"`
	EndAt      time.Time `json:"endAt"`
}

type Inventory struct {
	PlayerID int `json:"player_id"`
	ItemID   int `json:"item_id"`
	Quantity int `json:"quantity"`
}

type Items struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        int    `json:"type"` // 0 = Equipable, 1 = Consomable, 2 = Quest, 3 = Resources

}

type ItemStat struct {
	ItemID           int
	Strength         int `json:"strength"`
	Constitution     int `json:"constitution"`
	Mana             int `json:"mana"`
	Stamina          int `json:"stamina"`
	Dexterity        int `json:"dexterity"`
	Intelligence     int `json:"intelligence"`
	Wisdom           int `json:"wisdom"`
	Charisma         int `json:"charisma"`
	EnchantmentLevel int `json:"enchantmentLevel"`
}

type ItemEmplacement struct {
	ItemID      int    `json:"itemID"`
	Emplacement string `json:"emplacement"`
}

type ItemComplete struct {
	Item  Items
	Stats ItemStat
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

type Locations struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	IsSafety   bool   `json:"isSafety"`
	Difficulty int    `json:"difficulty"`
	Type       int    `json:"type"` // 0 = Resource, 1 = Mobs, 2 = City
	Size       int    `json:"size"`
}

type Loot struct {
	ID          int
	MobID       int
	ItemID      int
	QuantityMax int
	Rarity      int
}

type MobSkill struct {
	MobId   int
	SkillID int
}

type MobSpawns struct {
	MobID         int `json:"creatureID"`
	LocationID    int `json:"emplacementID"`
	LevelRequired int `json:"levelRequired"`
	SpawnRate     int `json:"spawnRate"`
}

type Mob struct {
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

type PetsMounts struct {
	ID          int  `json:"id"`
	MobID       int  `json:"creature_id"`
	IsMountable bool `json:"is_mountable"`
	Speed       int  `json:"speed"` // 0 = slow, 1 = normal, 2 = fast
}

type PlayerJobSkill struct {
	PlayerId int
	SkillID  int
}

type PlayerPet struct {
	PetID    int `json:"pet_id"`
	PlayerId int `json:"user_id"`
}

type PlayerSkill struct {
	PlayerId int
	SkillID  int
}

type Player struct {
	ID            int    `json:"id"`
	UserID        int    `json:"userID"`
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

type PlayerAction struct {
	PlayerId  int       `json:"user_id"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
	EndAt     time.Time `json:"end_at"`
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

type Race struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"` // description of the race
	Mana        int    `json:"mana"`
	Stamina     int    `json:"stamina"`
	Wisdom      int    `json:"wisdom"`
	Charisma    int    `json:"charisma"`
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

type Skill struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"` // physic or magic
	Description string `json:"description"`
}

type SkillStat struct {
	SkillID      int
	Strength     int `json:"strength"`
	Constitution int `json:"constitution"`
	Mana         int `json:"mana"`
	Stamina      int `json:"stamina"`
	Dexterity    int `json:"dexterity"`
	Intelligence int `json:"intelligence"`
	Wisdom       int `json:"wisdom"`
	Charisma     int `json:"charisma"`
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
