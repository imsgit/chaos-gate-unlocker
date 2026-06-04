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
	var assassinsAvailable, advancedTime bool

	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.GameUnlocksSaveState:
			object := record.SerializedObject.(*objects.GameUnlocksSaveState)
			for i := range object.Unlocks {
				if object.Unlocks[i].ID == TaintedSonsActivated {
					assassinsAvailable = true
				}
				if object.Unlocks[i].ID == KoramarMissionDefeated {
					advancedTime = true
				}
				if object.Unlocks[i].ID == AssassinsUnlocked {
					return false, true
				}
			}
		case internal.TimelineEventOccasion:
			object := record.SerializedObject.(*objects.TimelineEventOccasion)
			if object.EventToPlay.Key == ExecutionForce {
				return object.TriggerTime > 0, object.TriggerTime == 0
			}
		}
	}

	return assassinsAvailable && advancedTime, false
}
