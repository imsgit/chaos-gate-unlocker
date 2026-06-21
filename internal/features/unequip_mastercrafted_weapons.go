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
					weapon.Key = TechmarinePrefix + stem(strings.TrimPrefix(weapon.Key, TechmarinePrefix))
				} else if strings.HasPrefix(weapon.Key, MarketingPrefix) {
					weapon.Key = stem(strings.TrimPrefix(weapon.Key, MarketingPrefix))
				} else {
					weapon.Key = stem(weapon.Key)
				}

				if renamed, ok := weaponRename[weapon.Key]; ok {
					weapon.Key = renamed
				}
			}
		case internal.DreadnoughtState:
			object := record.SerializedObject.(*objects.DreadnoughtState)
			object.EquippedWeapons[0].Key = "Dreadnought_DoomFist"
			object.EquippedWeapons[1].Key = "Dreadnought_Lascannon"
		case internal.CallidusAssassinState, internal.CulexusAssassinState, internal.EversorAssassinState, internal.VindicareAssassinState:
			object := record.SerializedObject.(*objects.AssassinState)
			for _, weapon := range object.EquippedWeapons {
				weapon.Key = stem(weapon.Key)
			}
		}
	}
}

func (m *Manager) CanUnequipMastercraftedWeapons() (bool, bool) {
	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.KnightState:
			object := record.SerializedObject.(*objects.KnightState)
			class := getClass(object.CurrentLevelData.Key)
			if class == GarranCrowClass {
				continue
			}

			for _, weapon := range object.EquippedWeapons {
				if mastercrafted(weapon.Key, MarketingPrefix, TechmarinePrefix) {
					return true, true
				}
			}
		case internal.DreadnoughtState:
			object := record.SerializedObject.(*objects.DreadnoughtState)
			for _, weapon := range object.EquippedWeapons {
				if mastercrafted(weapon.Key, DreadnoughtPrefix) {
					return true, true
				}
			}
		case internal.CallidusAssassinState, internal.CulexusAssassinState, internal.EversorAssassinState, internal.VindicareAssassinState:
			object := record.SerializedObject.(*objects.AssassinState)
			for _, weapon := range object.EquippedWeapons {
				if mastercrafted(weapon.Key) {
					return true, true
				}
			}
		}
	}

	return false, true
}
