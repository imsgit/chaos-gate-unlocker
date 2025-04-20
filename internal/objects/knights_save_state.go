package objects

import "github.com/goccy/go-json"

type KnightsSaveState struct {
	Knights                       []IntValue      `json:"knights"`
	Squad                         json.RawMessage `json:"squad"`
	SecondSquad                   json.RawMessage `json:"secondSquad"`
	Remains                       []IntValue      `json:"remains"`
	DaysUntilNextCroweStateChange int             `json:"daysUntilNextCroweStateChange"`
	HasBeenOfferedCrowe           bool            `json:"hasBeenOfferedCrowe"`
	PotentialKnightRecruit        json.RawMessage `json:"potentialKnightRecruit"`
	VindicareAssassin             json.RawMessage `json:"vindicareAssassin"`
	EversorAssassin               json.RawMessage `json:"eversorAssassin"`
	CulexusAssassin               json.RawMessage `json:"culexusAssassin"`
	CallidusAssassin              json.RawMessage `json:"callidusAssassin"`
}
