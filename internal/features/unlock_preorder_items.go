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
	object := first[objects.ArmorySaveState](m, internal.ArmourySaveState)
	if object == nil {
		return
	}
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
}

func (m *Manager) CanUnlockPreorderItems() bool {
	object := first[objects.ArmorySaveState](m, internal.ArmourySaveState)
	if object == nil {
		return true
	}
	for _, wargear := range object.UnlockedWargears {
		if wargear.Data.Key == DominaLiberDaemonica {
			return false
		}
	}

	return true
}
