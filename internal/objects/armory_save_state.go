package objects

type UnlockedItem struct {
	Upgrades []bool      `json:"upgrades,omitempty"`
	Data     StringValue `json:"data"`
}

type ArmorySaveState struct {
	UnlockedWeapons  []UnlockedItem `json:"unlockedWeapons"`
	UnlockedArmours  []UnlockedItem `json:"unlockedArmours"`
	UnlockedWargears []UnlockedItem `json:"unlockedWargears"`
}
