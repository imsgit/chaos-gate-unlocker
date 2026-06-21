package features

const (
	PuritySealsUnlocked = "Purity_Seals_Unlocked"
	PuritySeals         = "PuritySeals"
	PoxwalkerNecropsy   = "Poxwalker_Necropsy"
)

func (m *Manager) UnlockPuritySeals() {
	m.unlockTimelineEvent(PuritySeals, 2)
}

func (m *Manager) CanUnlockPuritySeals() (bool, bool) {
	return m.canUnlockTimelineEvent(timelineUnlock{
		eventKey:   PuritySeals,
		prereqID:   PoxwalkerNecropsy,
		unlockedID: PuritySealsUnlocked,
	})
}
