package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"
)

const (
	ExecutionForce          = "ExecutionForce"
	TaintedSonsActivated    = "TaintedSons_Activated"
	AssassinsUnlocked       = "Assassins_Unlocked"
	HasQueuedExecutionForce = "Has_Queued_Execution_Force"
)

func (m *Manager) UnlockAssassins() {
	forEach(m, internal.GameUnlocksSaveState, func(o *objects.GameUnlocksSaveState) {
		o.Unlocks = append(o.Unlocks, objects.Unlock{
			ID: HasQueuedExecutionForce,
		})
	})

	m.unlockTimelineEvent(ExecutionForce, 3)
}

func (m *Manager) CanUnlockAssassins() (bool, bool) {
	return m.canUnlockTimelineEvent(timelineUnlock{
		eventKey:       ExecutionForce,
		prereqID:       TaintedSonsActivated,
		unlockedID:     AssassinsUnlocked,
		requireKoramar: true,
	})
}
