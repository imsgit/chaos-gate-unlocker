package features

import (
	"chaos-gate-unlocker/internal/objects"
)

const (
	CroweAvailable                      = "CroweAvailable"
	HasSeenFirstGrandMasterReportOfAct2 = "HasSeenFirstGrandMasterReportOfAct2"
)

func (m *Manager) UnlockGarranCrowe() {
	forEach(m, func(o *objects.GameUnlocksSaveState) {
		o.Unlocks = append(o.Unlocks, objects.Unlock{
			ID: HasSeenFirstGrandMasterReportOfAct2,
		})
	})
	forEach(m, func(o *objects.KnightsSaveState) {
		o.DaysUntilNextCroweStateChange = 0
	})
}

func (m *Manager) CanUnlockGarranCrowe() (bool, bool) {
	var croweAvailable, advancedTime bool

	for _, record := range m.state.LinearRecords {
		switch object := record.SerializedObject.(type) {
		case *objects.GameUnlocksSaveState:
			croweAvailable = croweAvailable || hasUnlock(object, CroweAvailable)
			advancedTime = advancedTime || hasUnlock(object, KoramarMissionDefeated)
		case *objects.KnightsSaveState:
			if object.HasBeenOfferedCrowe || object.DaysUntilNextCroweStateChange == 0 {
				return false, true
			}
		}
	}

	return croweAvailable && advancedTime, false
}
