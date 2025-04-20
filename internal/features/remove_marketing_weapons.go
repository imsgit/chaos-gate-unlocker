package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"

	"strings"
)

var NoSpecialLaunchEquipment = "No_Special_Launch_Equipment"

func (m *Manager) RemoveMarketingWeapons() {
	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.GameUnlocksSaveState:
			object := record.SerializedObject.(*objects.GameUnlocksSaveState)
			object.Unlocks = append(object.Unlocks, objects.Unlock{
				ID: NoSpecialLaunchEquipment,
			})
		case internal.KnightState:
			object := record.SerializedObject.(*objects.KnightState)
			for _, weapon := range object.EquippedWeapons {
				if strings.HasPrefix(weapon.Key, MarketingPrefix) {
					weapon.Key = strings.Split(strings.TrimPrefix(weapon.Key, MarketingPrefix), "_")[0]
					switch weapon.Key {
					case "Sword":
						weapon.Key = "ForceSword"
					case "Hammer":
						weapon.Key = "DaemonHammer"
					}
				}
			}
		}
	}
}

func (m *Manager) CanRemoveMarketingWeapons() bool {
	for _, record := range m.state.LinearRecords {
		if record.TypeName == internal.GameUnlocksSaveState {
			object := record.SerializedObject.(*objects.GameUnlocksSaveState)
			for i := range object.Unlocks {
				if object.Unlocks[i].ID == NoSpecialLaunchEquipment {
					return false
				}
			}
			break
		}
	}

	return true
}
