package internal

import (
	"chaos-gate-unlocker/internal/objects"

	"strconv"

	"github.com/goccy/go-json"
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
)

var typeNameToObject = map[string]func() interface{}{
	CurrencySaveState:       func() interface{} { return &objects.CurrencySaveState{} },
	GameUnlocksSaveState:    func() interface{} { return &objects.GameUnlocksSaveState{} },
	StarMapMission:          func() interface{} { return &objects.StarMapMission{} },
	StarMapNodeModel:        func() interface{} { return &objects.StarMapNodeModel{} },
	ConstructionProject:     func() interface{} { return &objects.ConstructionProject{} },
	ResearchProject:         func() interface{} { return &objects.ResearchProject{} },
	KnightsSaveState:        func() interface{} { return &objects.KnightsSaveState{} },
	KnightState:             func() interface{} { return &objects.KnightState{} },
	CallidusAssassinState:   func() interface{} { return &objects.AssassinState{} },
	CulexusAssassinState:    func() interface{} { return &objects.AssassinState{} },
	EversorAssassinState:    func() interface{} { return &objects.AssassinState{} },
	VindicareAssassinState:  func() interface{} { return &objects.AssassinState{} },
	DreadnoughtState:        func() interface{} { return &objects.DreadnoughtState{} },
	TimeManagerSaveState:    func() interface{} { return &objects.TimeManagerSaveState{} },
	TimelineEventOccasion:   func() interface{} { return &objects.TimelineEventOccasion{} },
	ArmourySaveState:        func() interface{} { return &objects.ArmorySaveState{} },
	StarMapMissionSaveState: func() interface{} { return &objects.StarMapMissionSaveState{} },
}

func (r *LinearRecord) MarshalJSON() ([]byte, error) {
	serializedContents := r.SerializedContents

	if r.SerializedObject != nil {
		serializedObject, err := json.Marshal(r.SerializedObject)
		if err != nil {
			return nil, err
		}

		serializedContents = []byte(strconv.Quote(string(serializedObject)))
	}

	return json.Marshal(linearRecord{
		TypeName:           r.TypeName,
		AssetName:          r.AssetName,
		SerializedContents: serializedContents,
	})
}

func (r *LinearRecord) UnmarshalJSON(data []byte) error {
	var t linearRecord
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	r.TypeName = t.TypeName
	r.AssetName = t.AssetName

	newObject, exists := typeNameToObject[t.TypeName]
	if !exists {
		r.SerializedContents = t.SerializedContents
		return nil
	}

	unquoted, _ := strconv.Unquote(string(t.SerializedContents))

	r.SerializedObject = newObject()
	return json.Unmarshal([]byte(unquoted), r.SerializedObject)
}
