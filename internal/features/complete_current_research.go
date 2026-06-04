package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"
)

func (m *Manager) CompleteCurrentResearch() {
	forEach(m, internal.ResearchProject, func(o *objects.ResearchProject) {
		o.ResearchPointsLeft = 0
	})
}

func (m *Manager) CanCompleteCurrentResearch() (bool, bool) {
	researchPointsLeft := -1
	forEach(m, internal.ResearchProject, func(o *objects.ResearchProject) {
		researchPointsLeft = o.ResearchPointsLeft
	})

	return researchPointsLeft > 0, researchPointsLeft == 0
}
