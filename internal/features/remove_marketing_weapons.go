package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"

	"strings"
)

var NoSpecialLaunchEquipment = "No_Special_Launch_Equipment"

func (m *Manager) RemoveMarketingWeapons() {
	forEach(m, internal.GameUnlocksSaveState, func(o *objects.GameUnlocksSaveState) {
		o.Unlocks = append(o.Unlocks, objects.Unlock{
			ID: NoSpecialLaunchEquipment,
		})
	})
	forEach(m, internal.KnightState, func(o *objects.KnightState) {
		for _, weapon := range o.EquippedWeapons {
			if strings.HasPrefix(weapon.Key, MarketingPrefix) {
				weapon.Key = stem(strings.TrimPrefix(weapon.Key, MarketingPrefix))
				if renamed, ok := weaponRename[weapon.Key]; ok {
					weapon.Key = renamed
				}
			}
		}
	})
}

func (m *Manager) CanRemoveMarketingWeapons() bool {
	object := first[objects.GameUnlocksSaveState](m, internal.GameUnlocksSaveState)
	if object == nil {
		return true
	}
	for i := range object.Unlocks {
		if object.Unlocks[i].ID == NoSpecialLaunchEquipment {
			return false
		}
	}

	return true
}
