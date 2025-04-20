package objects

import (
	"github.com/goccy/go-json"
)

type StarMapMission struct {
	StoryMissionId                                  string          `json:"storyMissionId"`
	MissionId                                       int             `json:"missionId"`
	NodeId                                          int             `json:"nodeId"`
	CreationDay                                     int             `json:"creationDay"`
	SpawnDay                                        json.RawMessage `json:"spawnDay"`
	InvisibilityDuration                            int             `json:"invisibilityDuration"`
	BaseDuration                                    json.RawMessage `json:"baseDuration"`
	IsPermanentlyInvisible                          bool            `json:"isPermanentlyInvisible"`
	MissionStrain                                   int             `json:"missionStrain"`
	IsTechnophageMission                            bool            `json:"isTechnophageMission"`
	HasOverrideBloom                                bool            `json:"hasOverrideBloom"`
	OverrideBloom                                   json.RawMessage `json:"overrideBloom"`
	MapName                                         string          `json:"mapName"`
	SuccessConsequences                             json.RawMessage `json:"successConsequences"`
	FailureConsequences                             json.RawMessage `json:"failureConsequences"`
	CoreReward                                      json.RawMessage `json:"coreReward"`
	ResourceSlotReward                              json.RawMessage `json:"resourceSlotReward"`
	RangedWeaponReward                              string          `json:"rangedWeaponReward"`
	MeleeWeaponReward                               string          `json:"meleeWeaponReward"`
	ArmourReward                                    string          `json:"armourReward"`
	WargearReward                                   string          `json:"wargearReward"`
	NumBloomspawnChosenForGrowingOrSpreadingMission int             `json:"numBloomspawnChosenForGrowingOrSpreadingMission"`
	EnemyLeaderArchetypeIDs                         []string        `json:"enemyLeaderArchetypeIDs"`
	SpawnPlan                                       json.RawMessage `json:"spawnPlan"`
	OverrideCorruptionTitleKey                      string          `json:"overrideCorruptionTitleKey"`
	OverrideCorruptionDescriptionKey                string          `json:"overrideCorruptionDescriptionKey"`
	Deed                                            json.RawMessage `json:"deed"`
	FrigateChallengeRating                          float64         `json:"frigateChallengeRating"`
	DeedAccepted                                    bool            `json:"deedAccepted"`
	Revealed                                        bool            `json:"revealed"`
	ArrivalEvents                                   json.RawMessage `json:"arrivalEvents"`
	HasPlayedArrivalEvent                           bool            `json:"hasPlayedArrivalEvent"`
	PostCombatVictoryEvents                         json.RawMessage `json:"postCombatVictoryEvents"`
	PostCombatLossEvents                            json.RawMessage `json:"postCombatLossEvents"`
	RemoveMissionOnLeaveTeleportarium               bool            `json:"removeMissionOnLeaveTeleportarium"`
	TeleportariumReturnLocation                     int             `json:"teleportariumReturnLocation"`
}
