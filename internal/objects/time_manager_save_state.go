package objects

import "encoding/json"

type TimeManagerSaveState struct {
	CurrentTime      json.RawMessage `json:"currentTime"`
	CurrentIncidents json.RawMessage `json:"currentIncidents"`
	CurrentOccasions struct {
		Values []IntValue `json:"values"`
	} `json:"currentOccasions"`
}
