package objects

import "github.com/goccy/go-json"

type AssassinState struct {
	MenuCustomModelOverride string          `json:"menuCustomModelOverride"`
	CustomModelOverride     string          `json:"customModelOverride"`
	EquippedWeapons         []*StringValue  `json:"equippedWeapons"`
	StartingStatusEffects   json.RawMessage `json:"startingStatusEffects"`
	CustomDeathFlow         json.RawMessage `json:"customDeathFlow"`
	DisableAnimatorOverride bool            `json:"disableAnimatorOverride"`
	IsNotCustomizable       bool            `json:"isNotCustomizable"`
	CustomizationState      json.RawMessage `json:"customizationState"`
	SurnameOverride         string          `json:"surnameOverride"`
	SurnameIndex            int             `json:"surnameIndex"`
	GivenName               string          `json:"givenName"`
	GivenNameOverride       string          `json:"givenNameOverride"`
	CurrentSideMission      struct {
		MissionID string `json:"missionID"`
		DaysLeft  int    `json:"daysLeft"`
	} `json:"currentSideMission"`
	LostMaxHealth       int             `json:"lostMaxHealth"`
	ClassPerks          json.RawMessage `json:"classPerks"`
	CurrentXP           int             `json:"currentXP"`
	CurrentLevelData    StringValue     `json:"currentLevelData"`
	NextLevelData       json.RawMessage `json:"nextLevelData"`
	ArmourRef           StringValue     `json:"armourRef"`
	Talents             json.RawMessage `json:"talents"`
	CreationTimestamp   int64           `json:"creationTimestamp"`
	TempEquipmentRefs   json.RawMessage `json:"tempEquipmentRefs"`
	EquippedItemClasses []*StringValue  `json:"equippedItemClasses"`
	HealthState         struct {
		Status                  int     `json:"status"`
		RecoveryTimeLeft        float64 `json:"recoveryTimeLeft"`
		HealingSuspended        bool    `json:"healingSuspended"`
		PartialRecoveryTimeLeft float64 `json:"partialRecoveryTimeLeft"`
	} `json:"healthState"`
}
