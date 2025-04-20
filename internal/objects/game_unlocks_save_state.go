package objects

type GameUnlocksSaveState struct {
	Unlocks []Unlock `json:"unlocks"`
}

type Unlock struct {
	DayUnlocked int    `json:"dayUnlocked"`
	ID          string `json:"id"`
}
