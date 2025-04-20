package objects

type ConstructionProject struct {
	UnlockID       string `json:"_unlockID"`
	RepairUnlockID string `json:"_repairUnlockID"`
	DaysLeft       int    `json:"_daysLeft"`
}
