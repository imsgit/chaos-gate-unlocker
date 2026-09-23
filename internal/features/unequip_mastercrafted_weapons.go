package features

import (
	"chaos-gate-unlocker/internal/objects"

	"strings"
)

const (
	TechmarinePrefix  = "Techmarine_"
	MarketingPrefix   = "Marketing_"
	DreadnoughtPrefix = "Dreadnought_"
)

var weaponRename = map[string]string{
	"Sword":  "ForceSword",
	"Shield": "StormShield",
	"Hammer": "DaemonHammer",
}

func (m *Manager) UnequipMastercraftedWeapons() {
	for _, record := range m.state.LinearRecords {
		switch object := record.SerializedObject.(type) {
		case *objects.KnightState:
			if stem(object.CurrentLevelData.Key) == GarranCrowClass {
				continue
			}

			for _, weapon := range object.EquippedWeapons {
				if rest, ok := strings.CutPrefix(weapon.Key, TechmarinePrefix); ok {
					weapon.Key = TechmarinePrefix + stem(rest)
				} else if rest, ok := strings.CutPrefix(weapon.Key, MarketingPrefix); ok {
					weapon.Key = stem(rest)
				} else {
					weapon.Key = stem(weapon.Key)
				}

				if renamed, ok := weaponRename[weapon.Key]; ok {
					weapon.Key = renamed
				}
			}
		case *objects.DreadnoughtState:
			object.EquippedWeapons[0].Key = "Dreadnought_DoomFist"
			object.EquippedWeapons[1].Key = "Dreadnought_Lascannon"
		case *objects.AssassinState:
			for _, weapon := range object.EquippedWeapons {
				weapon.Key = stem(weapon.Key)
			}
		}
	}
}

func (m *Manager) CanUnequipMastercraftedWeapons() (bool, bool) {
	for _, record := range m.state.LinearRecords {
		switch object := record.SerializedObject.(type) {
		case *objects.KnightState:
			if stem(object.CurrentLevelData.Key) == GarranCrowClass {
				continue
			}
			for _, weapon := range object.EquippedWeapons {
				if mastercrafted(weapon.Key, MarketingPrefix, TechmarinePrefix) {
					return true, true
				}
			}
		case *objects.DreadnoughtState:
			for _, weapon := range object.EquippedWeapons {
				if mastercrafted(weapon.Key, DreadnoughtPrefix) {
					return true, true
				}
			}
		case *objects.AssassinState:
			for _, weapon := range object.EquippedWeapons {
				if mastercrafted(weapon.Key) {
					return true, true
				}
			}
		}
	}

	return false, true
}
