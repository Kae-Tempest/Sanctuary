create table races
(
    id          serial
        constraint races_pk
            primary key,
    name        varchar(30),
    description text,
    mana        integer,
    stamina     integer,
    wisdom      integer,
    charisma    integer
);

alter table races
    owner to postgres;

create table jobs
(
    id           serial
        constraint jobs_pk
            primary key,
    name         varchar(25),
    description  text,
    strength     integer,
    constitution integer,
    mana         integer,
    stamina      integer,
    dexterity    integer,
    intelligence integer,
    wisdom       integer,
    charisma     integer
);

alter table jobs
    owner to postgres;

create table items
(
    id          serial
        constraint items_pk
            primary key,
    name        varchar(255),
    description text,
    type        integer,
    rank        varchar
);

alter table items
    owner to postgres;

create table guilds
(
    id    integer default nextval('compagnies_id_seq'::regclass) not null
        constraint guild_pk
            primary key,
    name  varchar(32),
    owner integer
);

alter table guilds
    owner to postgres;

create table skills
(
    id          serial
        constraint skills_pk
            primary key,
    name        varchar(100),
    description text,
    type        varchar
);

alter table skills
    owner to postgres;

create table quests
(
    id          serial
        constraint quests_pk
            primary key,
    title       varchar(50),
    description text,
    is_group    boolean,
    difficulty  integer,
    objective   json,
    reward      json
);

alter table quests
    owner to postgres;

create table mobs
(
    id           integer default nextval('creatures_id_seq'::regclass) not null
        constraint creatures_pk
            primary key,
    name         varchar(50),
    is_pet       boolean,
    strength     integer,
    constitution integer,
    mana         integer,
    stamina      integer,
    dexterity    integer,
    intelligence integer,
    wisdom       integer,
    charisma     integer,
    level        integer,
    hp           integer
);

alter table mobs
    owner to postgres;

create table pets_mounts
(
    mob_id       integer
        constraint pets_mounts_mobs_id_fk
            references mobs,
    is_mountable boolean,
    speed        integer,
    id           serial
        constraint pets_mounts_pk
            primary key
);

alter table pets_mounts
    owner to postgres;

create table locations
(
    id         serial
        constraint locations_pk
            primary key,
    name       varchar(50),
    is_safety  boolean,
    difficulty integer,
    type       integer,
    size       integer
);

alter table locations
    owner to postgres;

create table characters
(
    id             integer default nextval('players_id_seq'::regclass) not null
        constraint characters_pk_2
            primary key,
    email          varchar(32)
        constraint characters_pk
            unique,
    username       varchar(32),
    race_id        integer
        constraint characters_races_id_fk
            references races,
    job_id         integer
        constraint characters_jobs_id_fk
            references jobs,
    exp            integer,
    level          smallint,
    guild_id       integer,
    inventory_size integer,
    po             bigint,
    location_id    integer
        constraint characters_locations_id_fk
            references locations,
    user_id        integer
);

alter table characters
    owner to postgres;

create table inventory
(
    character_id integer
        constraint inventory_players_id_fk
            references characters,
    item_id      integer
        constraint inventory_items_id_fk
            references items,
    quantity     integer
);

alter table inventory
    owner to postgres;

create table guilds_members
(
    id           serial
        constraint guilds_members_pk
            primary key,
    guilds_id    integer
        constraint guilds_members_guildss_id_fk
            references guilds,
    character_id integer
        constraint guilds_members_players_id_fk
            references characters
);

alter table guilds_members
    owner to postgres;

create table equipment
(
    character_id     integer
        constraint equipment_players_id_fk
            references characters,
    helmet           integer,
    chestplate       integer,
    leggings         integer,
    boots            integer,
    mainhand         integer,
    offhand          integer,
    accessory_slot_0 integer,
    accessory_slot_1 integer,
    accessory_slot_2 integer,
    accessory_slot_3 integer
);

alter table equipment
    owner to postgres;

create table stats
(
    character_id integer
        constraint stats_players_id_fk
            references characters,
    strength     integer,
    constitution integer,
    mana         integer,
    stamina      integer,
    dexterity    integer,
    intelligence integer,
    wisdom       integer,
    charisma     integer,
    hp           integer
);

alter table stats
    owner to postgres;

create table summons_beats
(
    id           serial
        constraint summons_beats_pk
            primary key,
    character_id integer
        constraint summons_beats_players_id_fk
            references characters,
    name         varchar(50),
    strength     integer,
    constitution integer,
    mana         integer,
    stamina      integer,
    dexterity    integer,
    intelligence integer,
    wisdom       integer,
    charisma     integer
);

alter table summons_beats
    owner to postgres;

create table character_pets_mounts
(
    pet_id       integer
        constraint user_pets_mounts_pets_mounts_id_fk
            references pets_mounts,
    character_id integer
        constraint user_pets_mounts_players_id_fk
            references characters
);

alter table character_pets_mounts
    owner to postgres;

create table character_actions
(
    character_id integer
        constraint players_actions_players_id_fk
            references characters,
    action       varchar(50),
    created_at   timestamp,
    end_at       timestamp
);

alter table character_actions
    owner to postgres;

create table resources
(
    id                 integer default nextval('ressources_id_seq'::regclass) not null
        constraint ressources_pk
            primary key,
    name               varchar(50),
    location_id        integer
        constraint ressources_locations_id_fk
            references locations,
    quantities_per_min integer,
    item_id            integer
        constraint resources_items_id_fk
            references items
);

alter table resources
    owner to postgres;

create table mob_spawn
(
    mob_id         integer
        constraint creature_spawn_mobs_id_fk
            references mobs,
    location_id    integer
        constraint creature_spawn_locations_id_fk
            references locations,
    level_required integer,
    spawn_rate     double precision
);

alter table mob_spawn
    owner to postgres;

create table character_skill
(
    character_id integer
        constraint user_skill_players_id_fk
            references characters,
    skill_id     integer
        constraint user_skill_skills_id_fk
            references skills
);

alter table character_skill
    owner to postgres;

create table mob_skill
(
    mob_id   integer
        constraint creature_skill_mobs_id_fk
            references mobs,
    skill_id integer
        constraint creature_skill_skills_id_fk
            references skills
);

alter table mob_skill
    owner to postgres;

create table hunt_action
(
    character_id integer
        constraint hunt_action_players_id_fk
            references characters,
    location_id  integer
        constraint hunt_action_locations_id_fk
            references locations,
    mob_id       integer
        constraint hunt_action_mobs_id_fk
            references mobs,
    start_at     timestamp,
    end_at       timestamp
);

alter table hunt_action
    owner to postgres;

create table item_stats
(
    item_id           integer
        constraint item_stats_items_id_fk
            references items,
    strength          integer,
    constitution      integer,
    mana              integer,
    stamina           integer,
    dexterity         integer,
    intelligence      integer,
    wisdom            integer,
    charisma          integer,
    enchantment_level integer
);

alter table item_stats
    owner to postgres;

create table skill_stats
(
    skill_id     integer
        constraint skill_stat_skills_id_fk
            references skills,
    strength     integer,
    constitution integer,
    mana         integer,
    stamina      integer,
    dexterity    integer,
    intelligence integer,
    wisdom       integer,
    charisma     integer
);

alter table skill_stats
    owner to postgres;

create table loots
(
    id           integer not null
        constraint loots_pk
            primary key,
    mob_id       integer
        constraint loots_mobs_id_fk
            references mobs,
    item_id      integer
        constraint loots_items_id_fk
            references items,
    quantity_max integer,
    rarity       integer
);

alter table loots
    owner to postgres;

create table item_emplacement
(
    item_id     integer
        constraint item_emplacement_items_id_fk
            references items,
    emplacement integer
);

alter table item_emplacement
    owner to postgres;

create table users
(
    email      varchar,
    password   varchar,
    created_at timestamp,
    updated_at timestamp
);

alter table users
    owner to postgres;

create table job_skill
(
    job_id       integer     not null
        constraint job_skill_jobs_id_fk
            references jobs,
    name         varchar(50) not null,
    type         varchar(6),
    description  text,
    strength     integer,
    constitution integer,
    mana         integer,
    stamina      integer,
    dexterity    integer,
    intelligence integer,
    wisdom       integer,
    charisma     integer,
    id           serial
        constraint job_skill_pk
            primary key
);

alter table job_skill
    owner to postgres;

create table character_job_skill
(
    character_id integer
        constraint user_job_skill_players_id_fk
            references characters,
    job_skill_id integer
        constraint user_job_skill_job_skill_id_fk
            references job_skill
);

alter table character_job_skill
    owner to postgres;

