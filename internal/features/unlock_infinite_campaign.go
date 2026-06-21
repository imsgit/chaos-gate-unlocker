package features

import (
	"strings"

	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"
)

const (
	InfiniteCampaignTriggerTime = 999999

	DeadlineOccasionMarker = "TOO_LONG"
)

func (m *Manager) UnlockInfiniteCampaign() {
	forEach(m, internal.LoseGameOccasion, func(o *objects.LoseGameOccasion) {
		if strings.Contains(o.OccasionKey, DeadlineOccasionMarker) {
			o.TriggerTime = InfiniteCampaignTriggerTime
		}
	})
}

func (m *Manager) CanUnlockInfiniteCampaign() (bool, bool) {
	var hasClock, canUnlock bool
	forEach(m, internal.LoseGameOccasion, func(o *objects.LoseGameOccasion) {
		if !strings.Contains(o.OccasionKey, DeadlineOccasionMarker) {
			return
		}
		hasClock = true
		if o.TriggerTime < InfiniteCampaignTriggerTime {
			canUnlock = true
		}
	})
	return canUnlock, hasClock && !canUnlock
}
