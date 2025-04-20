package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"
)

const (
	CorruptedVesselReturn       = "CorruptedVessel_Return"
	CorruptedVesselNewEquipment = "CorruptedVessel_NewEquipment"
	DutyEternalActivated        = "DutyEternal_Activated"
	FrigateCleansed             = "Frigate_Cleansed"
)

func (m *Manager) UnlockGladiusFrigate() {
	var eventOccasion *objects.TimelineEventOccasion
	var eventOccasion2 *objects.TimelineEventOccasion
	var saveState *objects.TimeManagerSaveState

	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.TimelineEventOccasion:
			object := record.SerializedObject.(*objects.TimelineEventOccasion)
			if object.EventToPlay.Key == CorruptedVesselReturn {
				eventOccasion = object
			}
			if object.EventToPlay.Key == CorruptedVesselNewEquipment {
				eventOccasion2 = object
			}
		case internal.TimeManagerSaveState:
			object := record.SerializedObject.(*objects.TimeManagerSaveState)
			saveState = object
		}
	}

	if eventOccasion != nil && eventOccasion2 != nil {
		eventOccasion.TriggerTime = 0
		eventOccasion.SavedChosenResults.Values = []interface{}{}
		eventOccasion2.TriggerTime = 0
		eventOccasion2.SavedChosenResults.Values = []interface{}{}
		return
	}

	if saveState == nil {
		return
	}

	id := m.generateNewInstanceId()

	eventOccasion = &objects.TimelineEventOccasion{}
	eventOccasion.EventToPlay.Key = CorruptedVesselReturn
	eventOccasion.CalendarType = 15
	eventOccasion.SavedChosenResults.Values = []interface{}{}

	saveState.CurrentOccasions.Values = append(saveState.CurrentOccasions.Values, objects.IntValue{Key: id})
	m.state.LinearInstanceIds = append(m.state.LinearInstanceIds, id)
	m.state.LinearRecords = append(m.state.LinearRecords, &internal.LinearRecord{
		TypeName:         internal.TimelineEventOccasion,
		SerializedObject: eventOccasion,
	})
}

func (m *Manager) CanUnlockGladiusFrigate() (bool, bool) {
	var frigateAvailable, advancedTime bool

	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.GameUnlocksSaveState:
			object := record.SerializedObject.(*objects.GameUnlocksSaveState)
			for i := range object.Unlocks {
				if object.Unlocks[i].ID == DutyEternalActivated {
					frigateAvailable = true
				}
				if object.Unlocks[i].ID == KoramarMissionDefeated {
					advancedTime = true
				}
				if object.Unlocks[i].ID == FrigateCleansed {
					return false, true
				}
			}
		case internal.TimelineEventOccasion:
			object := record.SerializedObject.(*objects.TimelineEventOccasion)
			if object.EventToPlay.Key == CorruptedVesselReturn {
				return object.TriggerTime > 0, object.TriggerTime == 0
			}
		}
	}

	return frigateAvailable && advancedTime, false
}
