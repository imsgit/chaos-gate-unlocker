package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"

	"strings"
)

const (
	SynskinBodyglovePrefix = "SynskinBodyglove_"
)

var armorWithIncreasedSlots = map[string][]int{
	"PowerArmour_Plague_2":               {0, 3},
	"PowerArmour_HP_2":                   {0, 3},
	"PowerArmour_Armour_2":               {0, 2},
	"PowerArmour_HP_3":                   {0, 2},
	"PowerArmour_HP_Resist_3":            {0, 2},
	"PowerArmour_Armour_3":               {0, 2},
	"TerminatorArmour_HP_3":              {0, 2},
	"TechMarineArmour_WillPower_1":       {0, 3},
	"TechMarineArmour_Autos_1":           {1, 0},
	"TechMarineArmour_Vengeance_2":       {1, 2},
	"TechMarineArmour_Armour_2":          {0, 3},
	"TechMarineArmour_Autos_3":           {1, 4},
	"TechMarineArmour_Aegis_3":           {1, 2},
	"SynskinBodyglove_Eversor_Passive_1": {1, 0},
	"SynskinBodyglove_Eversor_HP_2":      {1, 0},
	"SynskinBodyglove_Callidus_Skull_1":  {1, 1},
	"SynskinBodyglove_Callidus_Skull_3":  {1, 1},
	"SynskinBodyglove_Callidus_Armour_3": {0, 2},
}

func (m *Manager) UnequipMastercraftedArmor() {
	upgrades := map[string][]bool{}

	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.ArmourySaveState:
			object := record.SerializedObject.(*objects.ArmorySaveState)
			for _, armor := range object.UnlockedArmours {
				if _, ok := armorWithIncreasedSlots[armor.Data.Key]; ok {
					upgrades[armor.Data.Key] = armor.Upgrades
				}
			}
		case internal.KnightState:
			object := record.SerializedObject.(*objects.KnightState)
			class := getClass(object.CurrentLevelData.Key)
			if class == GarranCrowClass {
				continue
			}

			if armor, ok := armorWithIncreasedSlots[object.ArmourRef.Key]; ok {
				dec := armor[0]
				if armor[1] > 0 && upgrades[object.ArmourRef.Key][armor[1]] {
					dec++
				}
				for i := len(object.EquippedItemClasses) - 1; i > 0 && dec > 0; i-- {
					object.EquippedItemClasses[i].Key = ""
					dec--
				}
			}
			object.ArmourRef.Key = strings.Split(object.ArmourRef.Key, "_")[0]
		case internal.DreadnoughtState:
			object := record.SerializedObject.(*objects.DreadnoughtState)
			object.ArmourRef.Key = strings.Split(object.ArmourRef.Key, "_")[0]
		case internal.CallidusAssassinState, internal.CulexusAssassinState, internal.EversorAssassinState, internal.VindicareAssassinState:
			object := record.SerializedObject.(*objects.AssassinState)
			if strings.HasPrefix(object.ArmourRef.Key, SynskinBodyglovePrefix) {
				if armor, ok := armorWithIncreasedSlots[object.ArmourRef.Key]; ok {
					dec := armor[0]
					if armor[1] > 0 && upgrades[object.ArmourRef.Key][armor[1]] {
						dec++
					}
					for i := len(object.EquippedItemClasses) - 1; i > 0 && dec > 0; i-- {
						object.EquippedItemClasses[i].Key = ""
						dec--
					}
				}
				object.ArmourRef.Key = SynskinBodyglovePrefix + strings.Split(strings.TrimPrefix(object.ArmourRef.Key, SynskinBodyglovePrefix), "_")[0]
			} else {
				object.ArmourRef.Key = strings.Split(object.ArmourRef.Key, "_")[0]
			}
		}
	}
}

func (m *Manager) CanUnequipMastercraftedArmor() bool {
	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.KnightState:
			object := record.SerializedObject.(*objects.KnightState)
			class := getClass(object.CurrentLevelData.Key)
			if class != GarranCrowClass && strings.Contains(object.ArmourRef.Key, "_") {
				return true
			}
		case internal.DreadnoughtState:
			object := record.SerializedObject.(*objects.DreadnoughtState)
			if strings.Contains(object.ArmourRef.Key, "_") {
				return true
			}
		case internal.CallidusAssassinState, internal.CulexusAssassinState, internal.EversorAssassinState, internal.VindicareAssassinState:
			object := record.SerializedObject.(*objects.AssassinState)
			if strings.Contains(strings.TrimPrefix(object.ArmourRef.Key, SynskinBodyglovePrefix), "_") {
				return true
			}
		}
	}

	return false
}
