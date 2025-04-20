package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"
)

var (
	HammerCrysylix3      = "Hammer_Crysylix_3"
	DominaLiberDaemonica = "DominaLiberDaemonica"
)

func (m *Manager) UnlockPreorderItems() {
	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.ArmourySaveState:
			object := record.SerializedObject.(*objects.ArmorySaveState)
			object.UnlockedWeapons = append(object.UnlockedWeapons, objects.UnlockedItem{
				Upgrades: []bool{false, false, false, false, false},
				Data: objects.StringValue{
					Key: HammerCrysylix3,
				},
			})
			object.UnlockedWargears = append(object.UnlockedWargears, objects.UnlockedItem{
				Data: objects.StringValue{
					Key: DominaLiberDaemonica,
				},
			})
			return
		}
	}
}

func (m *Manager) CanUnlockPreorderItems() bool {
	for _, record := range m.state.LinearRecords {
		if record.TypeName == internal.ArmourySaveState {
			object := record.SerializedObject.(*objects.ArmorySaveState)
			for _, wargear := range object.UnlockedWargears {
				if wargear.Data.Key == DominaLiberDaemonica {
					return false
				}
			}
			break
		}
	}

	return true
}
