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

func (r *LinearRecord) MarshalJSON() ([]byte, error) {
	serializedContents := r.SerializedContents

	if r.SerializedObject != nil {
		serializedObject, err := json.Marshal(r.SerializedObject)
		if err != nil {
			return nil, err
		}

		quoted := strconv.Quote(string(serializedObject))
		serializedContents = []byte(quoted)
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

	typeNameToObject := map[string]interface{}{
		CurrencySaveState:       &objects.CurrencySaveState{},
		GameUnlocksSaveState:    &objects.GameUnlocksSaveState{},
		StarMapMission:          &objects.StarMapMission{},
		StarMapNodeModel:        &objects.StarMapNodeModel{},
		ConstructionProject:     &objects.ConstructionProject{},
		ResearchProject:         &objects.ResearchProject{},
		KnightsSaveState:        &objects.KnightsSaveState{},
		KnightState:             &objects.KnightState{},
		CallidusAssassinState:   &objects.AssassinState{},
		CulexusAssassinState:    &objects.AssassinState{},
		EversorAssassinState:    &objects.AssassinState{},
		VindicareAssassinState:  &objects.AssassinState{},
		DreadnoughtState:        &objects.DreadnoughtState{},
		TimeManagerSaveState:    &objects.TimeManagerSaveState{},
		TimelineEventOccasion:   &objects.TimelineEventOccasion{},
		ArmourySaveState:        &objects.ArmorySaveState{},
		StarMapMissionSaveState: &objects.StarMapMissionSaveState{},
	}

	if obj, exists := typeNameToObject[t.TypeName]; exists {
		unquoted, _ := strconv.Unquote(string(t.SerializedContents))
		serializedContents := []byte(unquoted)

		r.SerializedObject = obj
		err = json.Unmarshal(serializedContents, r.SerializedObject)
		if err != nil {
			return err
		}

		return nil
	}

	r.SerializedContents = t.SerializedContents

	return nil
}
