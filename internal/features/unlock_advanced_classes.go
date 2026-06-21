package features

const (
	UnlockAdvancedClasses = "UnlockAdvancedClasses"
	AdvancedClasses       = "AdvancedClasses"
)

func (m *Manager) UnlockAdvancedClasses() {
	m.unlockTimelineEvent(AdvancedClasses, 2)
}

func (m *Manager) CanUnlockAdvancedClasses() (bool, bool) {
	return m.canUnlockTimelineEvent(timelineUnlock{
		eventKey:   AdvancedClasses,
		prereqID:   KoramarMissionDefeated,
		unlockedID: UnlockAdvancedClasses,
	})
}
