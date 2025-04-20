package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"

	"strings"
)

const (
	TechmarinePrefix  = "Techmarine_"
	MarketingPrefix   = "Marketing_"
	DreadnoughtPrefix = "Dreadnought_"
)

func (m *Manager) UnequipMastercraftedWeapons() {
	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.KnightState:
			object := record.SerializedObject.(*objects.KnightState)
			class := getClass(object.CurrentLevelData.Key)
			if class == GarranCrowClass {
				continue
			}

			for _, weapon := range object.EquippedWeapons {
				if strings.HasPrefix(weapon.Key, TechmarinePrefix) {
					weapon.Key = TechmarinePrefix + strings.Split(strings.TrimPrefix(weapon.Key, TechmarinePrefix), "_")[0]
				} else if strings.HasPrefix(weapon.Key, MarketingPrefix) {
					weapon.Key = strings.Split(strings.TrimPrefix(weapon.Key, MarketingPrefix), "_")[0]
				} else {
					weapon.Key = strings.Split(weapon.Key, "_")[0]
				}

				switch weapon.Key {
				case "Sword":
					weapon.Key = "ForceSword"
				case "Shield":
					weapon.Key = "StormShield"
				case "Hammer":
					weapon.Key = "DaemonHammer"
				}
			}
		case internal.DreadnoughtState:
			object := record.SerializedObject.(*objects.DreadnoughtState)
			object.EquippedWeapons[0].Key = "Dreadnought_DoomFist"
			object.EquippedWeapons[1].Key = "Dreadnought_Lascannon"
		case internal.CallidusAssassinState, internal.CulexusAssassinState, internal.EversorAssassinState, internal.VindicareAssassinState:
			object := record.SerializedObject.(*objects.AssassinState)
			for _, weapon := range object.EquippedWeapons {
				weapon.Key = strings.Split(weapon.Key, "_")[0]
			}
		}
	}
}

func (m *Manager) CanUnequipMastercraftedWeapons() bool {
	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.KnightState:
			object := record.SerializedObject.(*objects.KnightState)
			class := getClass(object.CurrentLevelData.Key)
			if class == GarranCrowClass {
				continue
			}

			for _, weapon := range object.EquippedWeapons {
				if strings.Contains(strings.TrimPrefix(strings.TrimPrefix(weapon.Key, MarketingPrefix), TechmarinePrefix), "_") {
					return true
				}
			}
		case internal.DreadnoughtState:
			object := record.SerializedObject.(*objects.DreadnoughtState)
			for _, weapon := range object.EquippedWeapons {
				if strings.Contains(strings.TrimPrefix(weapon.Key, DreadnoughtPrefix), "_") {
					return true
				}
			}
		case internal.CallidusAssassinState, internal.CulexusAssassinState, internal.EversorAssassinState, internal.VindicareAssassinState:
			object := record.SerializedObject.(*objects.AssassinState)
			for _, weapon := range object.EquippedWeapons {
				if strings.Contains(weapon.Key, "_") {
					return true
				}
			}
		}
	}

	return false
}
