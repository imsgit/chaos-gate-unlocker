package features

import (
	"chaos-gate-unlocker/internal/objects"

	"strings"
)

const (
	HonourOfTheAncientsComplete = "Honour_Of_The_Ancients_Complete"
	Flowering                   = "Flowering"
)

func (m *Manager) AuthorizeDreadnoughtMissions() {
	forEach(m, func(o *objects.StarMapMission) {
		if canBeTechnophageMission(o) {
			o.IsTechnophageMission = true
		}
	})
}

func (m *Manager) CanAuthorizeDreadnoughtMissions() (bool, bool) {
	var dreadnoughtAvailable, hasAvailableMissions, hasMissions bool

	currentMissions := map[int]bool{}

	for i, record := range m.state.LinearRecords {
		switch object := record.SerializedObject.(type) {
		case *objects.GameUnlocksSaveState:
			if hasUnlock(object, HonourOfTheAncientsComplete) {
				dreadnoughtAvailable = true
			}
		case *objects.StarMapMissionSaveState:
			for _, mission := range object.CurrentMissions.Values {
				currentMissions[mission.Key] = true
			}
		case *objects.StarMapMission:
			if currentMissions[m.state.LinearInstanceIds[i]] {
				if object.IsTechnophageMission {
					hasMissions = true
				} else if canBeTechnophageMission(object) {
					hasAvailableMissions = true
				}
			}
		}
	}

	return dreadnoughtAvailable && hasAvailableMissions,
		dreadnoughtAvailable && !hasAvailableMissions && hasMissions
}

func canBeTechnophageMission(s *objects.StarMapMission) bool {
	return s.StoryMissionId == "" && !strings.HasSuffix(s.MapName, Flowering)
}
