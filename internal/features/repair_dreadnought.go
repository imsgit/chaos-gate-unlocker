package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"
)

func (m *Manager) RepairDreadnought() {
	object := first[objects.DreadnoughtState](m, internal.DreadnoughtState)
	if object == nil {
		return
	}
	object.LostMaxHealth = 0
	object.HealthState.Status = 0
	object.HealthState.HealingSuspended = false
	object.HealthState.RecoveryTimeLeft = 0
}

func (m *Manager) CanRepairDreadnought() (bool, bool) {
	object := first[objects.DreadnoughtState](m, internal.DreadnoughtState)
	if object == nil {
		return false, false
	}
	return object.HealthState.Status > 0, true
}
