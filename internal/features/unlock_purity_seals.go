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
	var eventOccasion *objects.TimelineEventOccasion
	var saveState *objects.TimeManagerSaveState

	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.TimelineEventOccasion:
			object := record.SerializedObject.(*objects.TimelineEventOccasion)
			if object.EventToPlay.Key == PuritySeals {
				eventOccasion = object
			}
		case internal.TimeManagerSaveState:
			object := record.SerializedObject.(*objects.TimeManagerSaveState)
			saveState = object
		}
	}

	if eventOccasion != nil {
		eventOccasion.TriggerTime = 0
		eventOccasion.SavedChosenResults.Values = []interface{}{}
		return
	}

	if saveState == nil {
		return
	}

	id := m.generateNewInstanceId()

	eventOccasion = &objects.TimelineEventOccasion{}
	eventOccasion.EventToPlay.Key = PuritySeals
	eventOccasion.CalendarType = 2
	eventOccasion.SavedChosenResults.Values = []interface{}{}

	saveState.CurrentOccasions.Values = append(saveState.CurrentOccasions.Values, objects.IntValue{Key: id})
	m.state.LinearInstanceIds = append(m.state.LinearInstanceIds, id)
	m.state.LinearRecords = append(m.state.LinearRecords, &internal.LinearRecord{
		TypeName:         internal.TimelineEventOccasion,
		SerializedObject: eventOccasion,
	})
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
