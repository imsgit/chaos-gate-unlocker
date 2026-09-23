package internal

import (
	"chaos-gate-unlocker/internal/objects"

	"encoding/json"
)

const (
	CurrencySaveState       = "GreyKnights.CurrencySaveState"
	GameUnlocksSaveState    = "GreyKnights.GameUnlocksSaveState"
	KnightsSaveState        = "GreyKnights.KnightsSaveState"
	ArmourySaveState        = "GreyKnights.ArmourySaveState"
	KnightState             = "GreyKnights.KnightState"
	ConstructionProject     = "GreyKnights.ConstructionProject"
	ResearchProject         = "GreyKnights.ResearchProject"
	DreadnoughtState        = "GreyKnights.DreadnoughtState"
	CallidusAssassinState   = "GreyKnights.CallidusAssassinState"
	CulexusAssassinState    = "GreyKnights.CulexusAssassinState"
	EversorAssassinState    = "GreyKnights.EversorAssassinState"
	VindicareAssassinState  = "GreyKnights.VindicareAssassinState"
	StarMapMission          = "GreyKnights.StarMapMission"
	StarMapNodeModel        = "GreyKnights.StarMapNodeModel"
	TimeManagerSaveState    = "GreyKnights.TimeManagerSaveState"
	TimelineEventOccasion   = "GreyKnights.TimelineEventOccasion"
	StarMapMissionSaveState = "GreyKnights.StarMapMissionSaveState"
	LoseGameOccasion        = "GreyKnights.LoseGameOccasion"
)

var typeNameToObject = map[string]func() any{
	CurrencySaveState:       func() any { return &objects.CurrencySaveState{} },
	GameUnlocksSaveState:    func() any { return &objects.GameUnlocksSaveState{} },
	StarMapMission:          func() any { return &objects.StarMapMission{} },
	StarMapNodeModel:        func() any { return &objects.StarMapNodeModel{} },
	ConstructionProject:     func() any { return &objects.ConstructionProject{} },
	ResearchProject:         func() any { return &objects.ResearchProject{} },
	KnightsSaveState:        func() any { return &objects.KnightsSaveState{} },
	KnightState:             func() any { return &objects.KnightState{} },
	CallidusAssassinState:   func() any { return &objects.AssassinState{} },
	CulexusAssassinState:    func() any { return &objects.AssassinState{} },
	EversorAssassinState:    func() any { return &objects.AssassinState{} },
	VindicareAssassinState:  func() any { return &objects.AssassinState{} },
	DreadnoughtState:        func() any { return &objects.DreadnoughtState{} },
	TimeManagerSaveState:    func() any { return &objects.TimeManagerSaveState{} },
	TimelineEventOccasion:   func() any { return &objects.TimelineEventOccasion{} },
	ArmourySaveState:        func() any { return &objects.ArmorySaveState{} },
	StarMapMissionSaveState: func() any { return &objects.StarMapMissionSaveState{} },
	LoseGameOccasion:        func() any { return &objects.LoseGameOccasion{} },
}

type linearRecord LinearRecord

func (r *LinearRecord) MarshalJSON() ([]byte, error) {
	out := linearRecord(*r)
	if r.SerializedObject != nil {
		serializedObject, err := json.Marshal(r.SerializedObject)
		if err != nil {
			return nil, err
		}
		if out.SerializedContents, err = json.Marshal(string(serializedObject)); err != nil {
			return nil, err
		}
	}
	return json.Marshal(out)
}

func (r *LinearRecord) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, (*linearRecord)(r)); err != nil {
		return err
	}

	newObject, exists := typeNameToObject[r.TypeName]
	if !exists {
		return nil
	}

	var unquoted string
	if err := json.Unmarshal(r.SerializedContents, &unquoted); err != nil || unquoted == "" {
		return nil
	}

	r.SerializedContents = nil
	r.SerializedObject = newObject()
	return json.Unmarshal([]byte(unquoted), r.SerializedObject)
}
