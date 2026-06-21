package features

const (
	CorruptedVesselReturn       = "CorruptedVessel_Return"
	CorruptedVesselNewEquipment = "CorruptedVessel_NewEquipment"
	DutyEternalActivated        = "DutyEternal_Activated"
	FrigateCleansed             = "Frigate_Cleansed"
)

func (m *Manager) UnlockGladiusFrigate() {
	m.unlockTimelineEvent(CorruptedVesselReturn, 15, CorruptedVesselNewEquipment)
}

func (m *Manager) CanUnlockGladiusFrigate() (bool, bool) {
	return m.canUnlockTimelineEvent(timelineUnlock{
		eventKey:       CorruptedVesselReturn,
		prereqID:       DutyEternalActivated,
		unlockedID:     FrigateCleansed,
		requireKoramar: true,
	})
}
