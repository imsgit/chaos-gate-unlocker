package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"
)

const (
	CroweAvailable                      = "CroweAvailable"
	HasSeenFirstGrandMasterReportOfAct2 = "HasSeenFirstGrandMasterReportOfAct2"
)

func (m *Manager) UnlockGarranCrowe() {
	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.GameUnlocksSaveState:
			object := record.SerializedObject.(*objects.GameUnlocksSaveState)
			object.Unlocks = append(object.Unlocks, objects.Unlock{
				ID: HasSeenFirstGrandMasterReportOfAct2,
			})
		case internal.KnightsSaveState:
			object := record.SerializedObject.(*objects.KnightsSaveState)
			object.DaysUntilNextCroweStateChange = 0
		}
	}
}

func (m *Manager) CanUnlockGarranCrowe() (bool, bool) {
	var croweAvailable, advancedTime bool

	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.GameUnlocksSaveState:
			object := record.SerializedObject.(*objects.GameUnlocksSaveState)
			for i := range object.Unlocks {
				if object.Unlocks[i].ID == CroweAvailable {
					croweAvailable = true
				}
				if object.Unlocks[i].ID == KoramarMissionDefeated {
					advancedTime = true
				}
			}
		case internal.KnightsSaveState:
			object := record.SerializedObject.(*objects.KnightsSaveState)
			if object.HasBeenOfferedCrowe || object.DaysUntilNextCroweStateChange == 0 {
				return false, true
			}
		}
	}

	return croweAvailable && advancedTime, false
}
