package objects

import "encoding/json"

type StarMapMissionSaveState struct {
	WeekThatKoramarWasDefeated                   int             `json:"weekThatKoramarWasDefeated"`
	BloomEruptionNumber                          json.RawMessage `json:"bloomEruptionNumber"`
	SkippedWeeks                                 int             `json:"skippedWeeks"`
	DaysRemainingToNextBloomEruption             int             `json:"daysRemainingToNextBloomEruption"`
	LastFourMissionBloomEruption                 int             `json:"lastFourMissionBloomEruption"`
	BloomEruptionsRemainingUntilPoxusReactivates int             `json:"bloomEruptionsRemainingUntilPoxusReactivates"`
	Act3StartDay                                 int             `json:"act3StartDay"`
	Act3BloomEruptionNumber                      int             `json:"act3BloomEruptionNumber"`
	DateThatThreatLevelStartedIncreasing         json.RawMessage `json:"dateThatThreatLevelStartedIncreasing"`
	CurrentMissions                              struct {
		Values []IntValue `json:"values"`
	} `json:"currentMissions"`
	ActiveBloomTypes                         json.RawMessage `json:"activeBloomTypes"`
	DoomTrack                                json.RawMessage `json:"doomTrack"`
	NumMissionsSinceLastSuccessfulDeed       int             `json:"numMissionsSinceLastSuccessfulDeed"`
	NumEruptionsSinceLastTentarusHiveMission int             `json:"numEruptionsSinceLastTentarusHiveMission"`
	AlreadyUsedMissions                      []string        `json:"alreadyUsedMissions"`
}
