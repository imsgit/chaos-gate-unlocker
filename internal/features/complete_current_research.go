package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"
)

func (m *Manager) CompleteCurrentResearch() {
	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.ResearchProject:
			object := record.SerializedObject.(*objects.ResearchProject)
			object.ResearchPointsLeft = 0
		}
	}
}

func (m *Manager) CanCompleteCurrentResearch() (bool, bool) {
	researchPointsLeft := -1

	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.ResearchProject:
			object := record.SerializedObject.(*objects.ResearchProject)
			researchPointsLeft = object.ResearchPointsLeft
		}
	}

	return researchPointsLeft > 0, researchPointsLeft == 0
}
