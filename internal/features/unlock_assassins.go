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
	var eventOccasion *objects.TimelineEventOccasion
	var saveState *objects.TimeManagerSaveState

	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.GameUnlocksSaveState:
			object := record.SerializedObject.(*objects.GameUnlocksSaveState)
			object.Unlocks = append(object.Unlocks, objects.Unlock{
				ID: HasQueuedExecutionForce,
			})
		case internal.TimelineEventOccasion:
			object := record.SerializedObject.(*objects.TimelineEventOccasion)
			if object.EventToPlay.Key == ExecutionForce {
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
	eventOccasion.EventToPlay.Key = ExecutionForce
	eventOccasion.CalendarType = 3
	eventOccasion.SavedChosenResults.Values = []interface{}{}

	saveState.CurrentOccasions.Values = append(saveState.CurrentOccasions.Values, objects.IntValue{Key: id})
	m.state.LinearInstanceIds = append(m.state.LinearInstanceIds, id)
	m.state.LinearRecords = append(m.state.LinearRecords, &internal.LinearRecord{
		TypeName:         internal.TimelineEventOccasion,
		SerializedObject: eventOccasion,
	})
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
