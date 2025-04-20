package objects

type TimelineEventOccasion struct {
	EventToPlay                    StringValue `json:"eventToPlay"`
	TriggerTime                    float64     `json:"triggerTime"`
	ResetRegularEventsWhenTriggerd bool        `json:"resetRegularEventsWhenTriggerd"`
	SavedChosenResults             struct {
		Values []interface{} `json:"values"`
	} `json:"savedChosenResults"`
	CalendarTitleKey       string `json:"calendarTitleKey"`
	CalendarDescriptionKey string `json:"calendarDescriptionKey"`
	CalendarType           int    `json:"calendarType"`
}
