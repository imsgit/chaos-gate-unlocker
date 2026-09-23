package features

import (
	"chaos-gate-unlocker/internal/objects"

	"slices"
)

var preservedPerks = map[string]bool{
	"Champion":                true,
	"ExperimentalInoculation": true,
	"Dreadnought_Belligerent": true,
	"Dreadnought_Disciplined": true,
}

func retrainable(unit any) (perks *[]*objects.StringValue, class string) {
	switch u := unit.(type) {
	case *objects.KnightState:
		return &u.ClassPerks, stem(u.CurrentLevelData.Key)
	case *objects.DreadnoughtState:
		return &u.ClassPerks, stem(u.CurrentLevelData.Key)
	default:
		return nil, ""
	}
}

func trainedPerk(class string) func(*objects.StringValue) bool {
	defaultPerk := class + "_DefaultPerk"
	return func(perk *objects.StringValue) bool {
		return perk.Key != defaultPerk && !preservedPerks[perk.Key]
	}
}

func (m *Manager) RetrainUnit(unit any) {
	if perks, class := retrainable(unit); perks != nil {
		*perks = slices.DeleteFunc(*perks, trainedPerk(class))
	}
}

func (m *Manager) CanRetrainUnit(unit any) (enable, show bool) {
	perks, class := retrainable(unit)
	if perks == nil || class == GarranCrowClass {
		return false, false
	}
	return slices.ContainsFunc(*perks, trainedPerk(class)), true
}
