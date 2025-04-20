package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"

	"strings"
)

const (
	HonourOfTheAncientsComplete = "Honour_Of_The_Ancients_Complete"
	Flowering                   = "Flowering"
)

func (m *Manager) AuthorizeDreadnoughtMissions() {
	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.StarMapMission:
			object := record.SerializedObject.(*objects.StarMapMission)
			if canBeTechnophageMission(object) {
				object.IsTechnophageMission = true
			}
		}
	}
}

func (m *Manager) CanAuthorizeDreadnoughtMissions() (bool, bool) {
	var dreadnoughtAvailable, hasAvailableMissions, hasMissions bool

	currentMissions := map[int]bool{}

	for i, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.GameUnlocksSaveState:
			object := record.SerializedObject.(*objects.GameUnlocksSaveState)
			for i := range object.Unlocks {
				if object.Unlocks[i].ID == HonourOfTheAncientsComplete {
					dreadnoughtAvailable = true
					break
				}
			}
		case internal.StarMapMissionSaveState:
			object := record.SerializedObject.(*objects.StarMapMissionSaveState)
			for _, mission := range object.CurrentMissions.Values {
				currentMissions[mission.Key] = true
			}
		case internal.StarMapMission:
			object := record.SerializedObject.(*objects.StarMapMission)
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
