package features

import (
	"chaos-gate-unlocker/internal/objects"
)

func (m *Manager) HealUnit(unit any) {
	switch object := unit.(type) {
	case *objects.KnightState:
		object.LostMaxHealth = 0
		object.HealthState.Status = 0
		object.HealthState.HealingSuspended = false
		object.HealthState.PartialRecoveryTimeLeft = 0
		object.HealthState.RecoveryTimeLeft = 0
	case *objects.AssassinState:
		object.LostMaxHealth = 0
		object.HealthState.Status = 0
		object.HealthState.HealingSuspended = false
		object.HealthState.PartialRecoveryTimeLeft = 0
		object.HealthState.RecoveryTimeLeft = 0
	}
}

func (m *Manager) CanHealUnit(unit any) (bool, bool) {
	switch object := unit.(type) {
	case *objects.KnightState:
		return object.HealthState.Status > 0, true
	case *objects.AssassinState:
		return object.HealthState.Status > 0, true
	case *objects.DreadnoughtState:
		return object.HealthState.Status > 0, false
	}

	return false, false
}
