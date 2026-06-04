package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"
)

func (m *Manager) CompleteCurrentConstruction() {
	forEach(m, internal.ConstructionProject, func(o *objects.ConstructionProject) {
		o.DaysLeft = 0
	})
}

func (m *Manager) CanCompleteCurrentConstruction() (bool, bool) {
	daysLeft := -1
	forEach(m, internal.ConstructionProject, func(o *objects.ConstructionProject) {
		daysLeft = o.DaysLeft
	})

	return daysLeft > 0, daysLeft == 0
}
