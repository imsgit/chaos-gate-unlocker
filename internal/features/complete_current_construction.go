package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"
)

func (m *Manager) CompleteCurrentConstruction() {
	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.ConstructionProject:
			object := record.SerializedObject.(*objects.ConstructionProject)
			object.DaysLeft = 0
		}
	}
}

func (m *Manager) CanCompleteCurrentConstruction() (bool, bool) {
	daysLeft := -1

	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.ConstructionProject:
			object := record.SerializedObject.(*objects.ConstructionProject)
			daysLeft = object.DaysLeft
		}
	}

	return daysLeft > 0, daysLeft == 0
}
