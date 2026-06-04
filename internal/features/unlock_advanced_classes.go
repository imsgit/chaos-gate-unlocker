package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"
)

const (
	UnlockAdvancedClasses = "UnlockAdvancedClasses"
	AdvancedClasses       = "AdvancedClasses"
)

func (m *Manager) UnlockAdvancedClasses() {
	m.unlockTimelineEvent(AdvancedClasses, 2)
}

func (m *Manager) CanUnlockAdvancedClasses() (bool, bool) {
	var advancedTime bool

	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.GameUnlocksSaveState:
			object := record.SerializedObject.(*objects.GameUnlocksSaveState)
			for i := range object.Unlocks {
				if object.Unlocks[i].ID == KoramarMissionDefeated {
					advancedTime = true
				}
				if object.Unlocks[i].ID == UnlockAdvancedClasses {
					return false, true
				}
			}
		case internal.TimelineEventOccasion:
			object := record.SerializedObject.(*objects.TimelineEventOccasion)
			if object.EventToPlay.Key == AdvancedClasses {
				return object.TriggerTime > 0, object.TriggerTime == 0
			}
		}
	}

	return advancedTime, false
}
