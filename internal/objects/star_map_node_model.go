package objects

import "github.com/goccy/go-json"

type StarMapNodeModel struct {
	NodeId                             int             `json:"nodeId"`
	Act                                int             `json:"act"`
	Position                           json.RawMessage `json:"position"`
	AllPossiblePlanetNameKeys          json.RawMessage `json:"allPossiblePlanetNameKeys"`
	PlanetNameKeyOverride              json.RawMessage `json:"planetNameKeyOverride"`
	Links                              json.RawMessage `json:"links"`
	Bloom                              json.RawMessage `json:"bloom"`
	NumTimesSkippedForFloweringMission int             `json:"numTimesSkippedForFloweringMission"`
	IsDestroyed                        json.RawMessage `json:"isDestroyed"`
	HasBeenExterminatused              json.RawMessage `json:"hasBeenExterminatused"`
	HasPrognosticar                    struct {
		Value bool `json:"value"`
	} `json:"hasPrognosticar"`
	IsHidden              json.RawMessage `json:"isHidden"`
	WarpStormsWeAreIn     json.RawMessage `json:"warpStormsWeAreIn"`
	MissionBackgroundName json.RawMessage `json:"missionBackgroundName"`
	OrbVisualIndex        int             `json:"orbVisualIndex"`
}
