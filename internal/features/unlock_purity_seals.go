package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"
)

const (
	PuritySealsUnlocked = "Purity_Seals_Unlocked"
	PuritySeals         = "PuritySeals"
	PoxwalkerNecropsy   = "Poxwalker_Necropsy"
)

func (m *Manager) UnlockPuritySeals() {
	m.unlockTimelineEvent(PuritySeals, 2)
}

func (m *Manager) CanUnlockPuritySeals() (bool, bool) {
	var seedsAvailable bool

	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.GameUnlocksSaveState:
			object := record.SerializedObject.(*objects.GameUnlocksSaveState)
			for i := range object.Unlocks {
				if object.Unlocks[i].ID == PoxwalkerNecropsy {
					seedsAvailable = true
				}
				if object.Unlocks[i].ID == PuritySealsUnlocked {
					return false, true
				}
			}
		case internal.TimelineEventOccasion:
			object := record.SerializedObject.(*objects.TimelineEventOccasion)
			if object.EventToPlay.Key == PuritySeals {
				return object.TriggerTime > 0, object.TriggerTime == 0
			}
		}
	}

	return seedsAvailable, false
}
