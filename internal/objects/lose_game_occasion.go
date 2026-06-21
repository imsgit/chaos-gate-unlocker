package objects

import "github.com/goccy/go-json"

type LoseGameOccasion struct {
	OccasionKey           string          `json:"occasionKey"`
	TriggerTime           float64         `json:"triggerTime"`
	GameLossCinematicPath json.RawMessage `json:"gameLossCinematicPath"`
	PopupTitle            json.RawMessage `json:"popupTitle"`
	PopupSubtitle         json.RawMessage `json:"popupSubtitle"`
	PopupDescription      json.RawMessage `json:"popupDescription"`
	PopupFooter           json.RawMessage `json:"popupFooter"`
	PopupTopper           json.RawMessage `json:"popupTopper"`
	FooterColour          json.RawMessage `json:"footerColour"`
}
