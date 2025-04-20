package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"
)

func (m *Manager) RepairDreadnought() {
	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.DreadnoughtState:
			object := record.SerializedObject.(*objects.DreadnoughtState)
			object.LostMaxHealth = 0
			object.HealthState.Status = 0
			object.HealthState.HealingSuspended = false
			object.HealthState.RecoveryTimeLeft = 0
			return
		}
	}
}

func (m *Manager) CanRepairDreadnought() (bool, bool) {
	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.DreadnoughtState:
			object := record.SerializedObject.(*objects.DreadnoughtState)
			return object.HealthState.Status > 0, true
		}
	}

	return false, false
}
