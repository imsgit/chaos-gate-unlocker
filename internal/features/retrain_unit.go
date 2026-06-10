package features

import (
	"chaos-gate-unlocker/internal/objects"
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
		return &u.ClassPerks, getClass(u.CurrentLevelData.Key)
	case *objects.DreadnoughtState:
		return &u.ClassPerks, getClass(u.CurrentLevelData.Key)
	default:
		return nil, ""
	}
}

func (m *Manager) RetrainUnit(unit any) {
	perks, class := retrainable(unit)
	if perks == nil {
		return
	}

	defaultPerk := class + "_DefaultPerk"

	kept := (*perks)[:0]
	for _, perk := range *perks {
		if perk.Key == defaultPerk || preservedPerks[perk.Key] {
			kept = append(kept, perk)
		}
	}
	*perks = kept
}

func (m *Manager) CanRetrainUnit(unit any) (enable, show bool) {
	perks, class := retrainable(unit)
	if perks == nil {
		return false, false
	}

	if class == GarranCrowClass {
		return false, false
	}

	defaultPerk := class + "_DefaultPerk"
	for _, perk := range *perks {
		if perk.Key != defaultPerk && !preservedPerks[perk.Key] {
			return true, true
		}
	}
	return false, true
}
