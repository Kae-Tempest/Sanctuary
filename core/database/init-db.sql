create sequence compagnies_id_seq
    as integer;

alter sequence compagnies_id_seq owner to postgres;

create sequence ressources_id_seq
    as integer;

alter sequence ressources_id_seq owner to postgres;

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
    id                serial
        constraint items_pk
            primary key,
    name              varchar(255),
    description       text,
    type              integer,
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

alter sequence compagnies_id_seq owned by guilds.id;

create table skills
(
    id           serial
        constraint skills_pk
            primary key,
    name         varchar(100),
    description  text,
    strength     integer,
    constitution integer,
    mana         integer,
    stamina      integer,
    dexterity    integer,
    intelligence integer,
    wisdom       integer,
    charisma     integer,
    type         varchar
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

create table creatures
(
    id           serial
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

alter table creatures
    owner to postgres;

create table pets_mounts
(
    creature_id  integer
        constraint pets_mounts_creatures_id_fk
            references creatures,
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
    is__safety boolean,
    difficulty integer,
    type       integer,
    size       integer
);

alter table locations
    owner to postgres;

create table players
(
    id             serial
        constraint players_pk
            primary key,
    email          varchar(32)
        constraint players_pk_2
            unique,
    username       varchar(32),
    race_id        integer
        constraint players_races_id_fk
            references races,
    job_id         integer
        constraint players_jobs_id_fk
            references jobs,
    exp            integer,
    level          smallint,
    guild_id       integer,
    inventory_size integer,
    po             bigint,
    location_id    integer
        constraint players_locations_id_fk
            references locations
);

alter table players
    owner to postgres;

create table inventory
(
    player_id integer
        constraint inventory_players_id_fk
            references players,
    item_id   integer
        constraint inventory_pk
            unique
        constraint inventory_items_id_fk
            references items,
    quantity  integer
);

alter table inventory
    owner to postgres;

create table guilds_members
(
    id        serial
        constraint guilds_members_pk
            primary key,
    guilds_id integer
        constraint guilds_members_guildss_id_fk
            references guilds,
    user_id   integer
        constraint guilds_members_players_id_fk
            references players
);

alter table guilds_members
    owner to postgres;

create table equipment
(
    player_id  integer
        constraint equipment_players_id_fk
            references players,
    helmet     integer,
    chestplate integer,
    leggings   integer,
    boots      integer,
    mainhand   integer,
    offhand    integer,
    accesory_0 integer,
    accesory_1 integer,
    accesory_2 integer,
    accesory_3 integer
);

alter table equipment
    owner to postgres;

create table stats
(
    player_id    integer
        constraint stats_players_id_fk
            references players,
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
    player_id    integer
        constraint summons_beats_players_id_fk
            references players,
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

create table user_pets_mounts
(
    pet_id    integer
        constraint user_pets_mounts_pets_mounts_id_fk
            references pets_mounts,
    player_id integer
        constraint user_pets_mounts_players_id_fk
            references players
);

alter table user_pets_mounts
    owner to postgres;

create table players_actions
(
    player_id  integer
        constraint players_actions_players_id_fk
            references players,
    action     varchar(50),
    created_at timestamp,
    end_at     timestamp
);

alter table players_actions
    owner to postgres;

create table resources_types
(
    id   serial
        constraint resources_types_pk
            primary key,
    name varchar(50)
);

alter table resources_types
    owner to postgres;

create table resources
(
    id                 integer default nextval('ressources_id_seq'::regclass) not null
        constraint ressources_pk
            primary key,
    name               varchar(50),
    emplacement_id     integer
        constraint ressources_locations_id_fk
            references locations,
    ressources_type_id integer
        constraint ressources_resources_types_id_fk
            references resources_types,
    quantities_per_min integer
);

alter table resources
    owner to postgres;

alter sequence ressources_id_seq owned by resources.id;

create table ressourceinventory
(
    user_id  integer
        constraint ressourceinventory_players_id_fk
            references players,
    item_id  integer
        constraint ressourceinventory_pk
            unique
        constraint ressourceinventory_resources_id_fk
            references resources,
    quantity integer
);

alter table ressourceinventory
    owner to postgres;

create table creaturespawn
(
    creature_id    integer
        constraint creaturespawn_creatures_id_fk
            references creatures,
    emplacement_id integer
        constraint creaturespawn_locations_id_fk
            references locations,
    level_required integer,
    spawn_rate     double precision
);

alter table creaturespawn
    owner to postgres;

create table job_skill
(
    id           integer     not null
        constraint job_skill_pk
            primary key,
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
    charisma     integer
);

alter table job_skill
    owner to postgres;

create table user_skill
(
    player_id integer
        constraint user_skill_players_id_fk
            references players,
    skill_id  integer
        constraint user_skill_skills_id_fk
            references skills
);

alter table user_skill
    owner to postgres;

create table user_job_skill
(
    player_id    integer
        constraint user_job_skill_players_id_fk
            references players,
    job_skill_id integer
        constraint user_job_skill_job_skill_id_fk
            references job_skill
);

alter table user_job_skill
    owner to postgres;

create table creature_skill
(
    creature_id integer
        constraint creature_skill_creatures_id_fk
            references creatures,
    skill_id    integer
        constraint creature_skill_skills_id_fk
            references skills
);

alter table creature_skill
    owner to postgres;

create table hunt_action
(
    player_id  integer
        constraint hunt_action_players_id_fk
            references players,
    btn_id     varchar(20),
    message_id varchar,
    channel_id varchar
);

alter table hunt_action
    owner to postgres;

