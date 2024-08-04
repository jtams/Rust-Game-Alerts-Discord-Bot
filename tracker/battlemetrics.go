package tracker

import "time"

// This files contains the struct data for the BattleMetrics API response
type BattleMetricsResponse struct {
	Data     Data     `json:"data"`
	Included []Player `json:"included"`
}

type Data struct {
	Type          string        `json:"type"`
	ID            string        `json:"id"`
	Attributes    Attributes    `json:"attributes"`
	Relationships Relationships `json:"relationships"`
}

type Attributes struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	IP          string    `json:"ip"`
	Port        int       `json:"port"`
	Players     int       `json:"players"`
	MaxPlayers  int       `json:"maxPlayers"`
	Rank        int       `json:"rank"`
	Location    []float64 `json:"location"`
	Status      string    `json:"status"`
	Details     Details   `json:"details"`
	Private     bool      `json:"private"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	PortQuery   int       `json:"portQuery"`
	Country     string    `json:"country"`
	QueryStatus string    `json:"queryStatus"`
}

type Details struct {
	Tags               []string      `json:"tags"`
	Official           bool          `json:"official"`
	RustType           string        `json:"rust_type"`
	Map                string        `json:"map"`
	Environment        string        `json:"environment"`
	RustBuild          string        `json:"rust_build"`
	RustEntCntI        int           `json:"rust_ent_cnt_i"`
	RustFPS            float32       `json:"rust_fps"`
	RustFPSAvg         float32       `json:"rust_fps_avg"`
	RustGCCl           int           `json:"rust_gc_cl"`
	RustGCMb           int           `json:"rust_gc_mb"`
	RustHash           string        `json:"rust_hash"`
	RustHeaderImage    string        `json:"rust_headerimage"`
	RustMemPv          interface{}   `json:"rust_mem_pv"`
	RustMemWs          interface{}   `json:"rust_mem_ws"`
	Pve                bool          `json:"pve"`
	RustUptime         int           `json:"rust_uptime"`
	RustURL            string        `json:"rust_url"`
	RustWorldSeed      int           `json:"rust_world_seed"`
	RustWorldSize      int           `json:"rust_world_size"`
	RustWorldLevelURL  string        `json:"rust_world_levelurl"`
	RustMaps           RustMaps      `json:"rust_maps"`
	RustDescription    string        `json:"rust_description"`
	RustModded         bool          `json:"rust_modded"`
	RustQueuedPlayers  int           `json:"rust_queued_players"`
	RustGamemode       string        `json:"rust_gamemode"`
	RustBorn           time.Time     `json:"rust_born"`
	RustLastEntDrop    time.Time     `json:"rust_last_ent_drop"`
	RustLastSeedChange time.Time     `json:"rust_last_seed_change"`
	RustLastWipe       time.Time     `json:"rust_last_wipe"`
	RustLastWipeEnt    int           `json:"rust_last_wipe_ent"`
	RustSettingsSource string        `json:"rust_settings_source"`
	RustSettings       RustSettings  `json:"rust_settings"`
	RustWipes          []interface{} `json:"rust_wipes"`
	ServerSteamID      string        `json:"serverSteamId"`
}

type RustMaps struct {
	Seed             int              `json:"seed"`
	Size             int              `json:"size"`
	URL              string           `json:"url"`
	ThumbnailURL     string           `json:"thumbnailUrl"`
	MonumentCount    int              `json:"monumentCount"`
	Barren           bool             `json:"barren"`
	UpdatedAt        time.Time        `json:"updatedAt"`
	MapURL           string           `json:"mapUrl"`
	BiomePercentages BiomePercentages `json:"biomePercentages"`
	Islands          int              `json:"islands"`
	Mountains        int              `json:"mountains"`
	IceLakes         int              `json:"iceLakes"`
	Rivers           int              `json:"rivers"`
	MonumentCounts   MonumentCounts   `json:"monumentCounts"`
	Monuments        []string         `json:"monuments"`
}

type BiomePercentages struct {
	S float64 `json:"s"`
	D float64 `json:"d"`
	F float64 `json:"f"`
	T float64 `json:"t"`
}

type MonumentCounts struct {
	FerryTerminal            int `json:"Ferry Terminal"`
	LargeHarbor              int `json:"Large Harbor"`
	SmallHarbor              int `json:"Small Harbor"`
	FishingVillage           int `json:"Fishing Village"`
	MilitaryBase             int `json:"Military Base"`
	ArcticResearchBase       int `json:"Arctic Research Base"`
	LaunchSite               int `json:"Launch Site"`
	Excavator                int `json:"Excavator"`
	Airfield                 int `json:"Airfield"`
	WaterTreatment           int `json:"Water Treatment"`
	Powerplant               int `json:"Powerplant"`
	Trainyard                int `json:"Trainyard"`
	MilitaryTunnels          int `json:"Military Tunnels"`
	NuclearMissileSilo       int `json:"Nuclear Missile Silo"`
	Outpost                  int `json:"Outpost"`
	SatelliteDish            int `json:"Satellite Dish"`
	SphereTank               int `json:"Sphere Tank"`
	Junkyard                 int `json:"Junkyard"`
	Ranch                    int `json:"Ranch"`
	SewerBranch              int `json:"Sewer Branch"`
	SulfurQuarry             int `json:"Sulfur Quarry"`
	LargeBarn                int `json:"Large Barn"`
	StoneQuarry              int `json:"Stone Quarry"`
	HqmQuarry                int `json:"Hqm Quarry"`
	TunnelEntrance           int `json:"Tunnel Entrance"`
	TunnelEntranceTransition int `json:"Tunnel Entrance Transition"`
	IceLake                  int `json:"Ice Lake"`
	Warehouse                int `json:"Warehouse"`
	Supermarket              int `json:"Supermarket"`
	GasStation               int `json:"Gas Station"`
	PowerSubstationSmall     int `json:"Power Substation Small"`
	PowerSubstationBig       int `json:"Power Substation Big"`
	Powerline                int `json:"Powerline"`
	CaveSmallEasy            int `json:"Cave Small Easy"`
	CaveSmallMedium          int `json:"Cave Small Medium"`
	CaveLargeSewersHard      int `json:"Cave Large Sewers Hard"`
	CaveSmallHard            int `json:"Cave Small Hard"`
	CaveMediumMedium         int `json:"Cave Medium Medium"`
	UnderwaterLab            int `json:"Underwater Lab"`
	LargeOilrig              int `json:"Large Oilrig"`
	SmallOilrig              int `json:"Small Oilrig"`
	Lighthouse               int `json:"Lighthouse"`
	Iceberg                  int `json:"Iceberg"`
}

type RustSettings struct {
	Upkeep        int           `json:"upkeep"`
	Blueprints    bool          `json:"blueprints"`
	ForceWipeType string        `json:"forceWipeType"`
	GroupLimit    int           `json:"groupLimit"`
	TeamUILimit   int           `json:"teamUILimit"`
	Kits          bool          `json:"kits"`
	Rates         Rates         `json:"rates"`
	Wipes         []interface{} `json:"wipes"`
	Decay         int           `json:"decay"`
	TimeZone      string        `json:"timeZone"`
	Version       int           `json:"version"`
}

type Rates struct {
	Component int `json:"component"`
	Craft     int `json:"craft"`
	Gather    int `json:"gather"`
	Scrap     int `json:"scrap"`
}

type Relationships struct {
	Game Game `json:"game"`
}

type Game struct {
	Data GameData `json:"data"`
}

type GameData struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type Player struct {
	Type          string              `json:"type"`
	ID            string              `json:"id"`
	Attributes    PlayerAttributes    `json:"attributes"`
	Relationships PlayerRelationships `json:"relationships"`
	Meta          Meta                `json:"meta"`
}

type PlayerAttributes struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Private       bool      `json:"private"`
	PositiveMatch bool      `json:"positiveMatch"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type PlayerRelationships struct {
	Server GameData `json:"server"`
}

type Meta struct {
	Metadata []Metadata `json:"metadata"`
}

type Metadata struct {
	Key     string      `json:"key"`
	Value   interface{} `json:"value"`
	Private bool        `json:"private"`
}
