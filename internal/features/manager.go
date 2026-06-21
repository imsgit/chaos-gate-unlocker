package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"

	"sort"
	"strconv"
	"strings"
)

const (
	GarranCrowClass  = "GarranCrowe"
	TechmarineClass  = "TechMarine"
	DreadnoughtClass = "Dreadnought"

	KoramarMissionDefeated = "Koramar_Mission_Defeated"
)

var Surnames = []string{
	"Adantor", "Aelvar", "Aldar", "Arelis", "Bhask", "Bors", "Cadulon", "Corvane", "Crassus", "Dalmar",
	"Decran", "Durant", "Dvorn", "Edeon", "Elgon", "Esdrios", "Garedian", "Garr", "Gul", "Hale",
	"Harne", "Ignatius", "Invio", "Iolanthus", "Issad", "Kai", "Kain", "Kalmar", "Kern", "Malchus",
	"Massius", "Myr", "Nedth", "Neodan", "Palamedes", "Pardum", "Phoros", "Rao", "Rugan", "Rythvane",
	"Santor", "Solor", "Sorak", "Storm", "Tarn", "Tekios", "Thawn", "Thule", "Tor", "Trevan",
	"Tydes", "Valdar", "Varn", "Vorn", "Vortimer", "Zaebus",
}

var AssassinSurnames = []string{
	"Asch", "Fang", "Five", "Ganus", "Garmeaux", "Koln", "Kord", "Lasc", "Nine", "Novac",
	"Pec", "Raithe", "Rhasc", "Skult", "Torq", "Vald", "Vanus", "Zhau",
}

type Manager struct {
	state *internal.State
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) ApplyState() func(state *internal.State) {
	return func(state *internal.State) {
		m.state = state
	}
}

type Sort []interface{}

func (s Sort) Len() int {
	return len(s)
}

func classStatusLvlName(obj interface{}) (class, status, lvl int, name string) {
	switch object := obj.(type) {
	case *objects.KnightState:
		status = object.HealthState.Status
		if object.CurrentSideMission.MissionID != "" {
			status = 5
		}
		return 2, status, getLvl(object.CurrentLevelData.Key), object.GivenName
	case *objects.AssassinState:
		status = object.HealthState.Status
		if object.CurrentSideMission.MissionID != "" {
			status = 5
		}
		return 1, 0, getLvl(object.CurrentLevelData.Key), object.GivenName
	case *objects.DreadnoughtState:
		return 0, 0, getLvl(object.CurrentLevelData.Key), object.GivenName
	default:
		return 0, 0, 0, ""
	}
}

func (s Sort) Less(i, j int) bool {
	iClass, iStatus, iLvl, iName := classStatusLvlName(s[i])
	jClass, jStatus, jLvl, jName := classStatusLvlName(s[j])

	if jClass != iClass {
		return jClass < iClass
	}
	if jStatus != iStatus {
		return jStatus < iStatus
	}
	if jLvl != iLvl {
		return jLvl < iLvl
	}
	return jName > iName
}

func (s Sort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (m *Manager) Units() []interface{} {
	var units []interface{}
	knightInBarracks := map[int]bool{}

	for i, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.KnightsSaveState:
			object := record.SerializedObject.(*objects.KnightsSaveState)
			for _, knight := range object.Knights {
				knightInBarracks[knight.Key] = true
			}
		case internal.KnightState:
			object := record.SerializedObject.(*objects.KnightState)
			if knightInBarracks[m.state.LinearInstanceIds[i]] {
				units = append(units, object)
			}
		case internal.DreadnoughtState:
			object := record.SerializedObject.(*objects.DreadnoughtState)
			units = append(units, object)
		case internal.CallidusAssassinState, internal.CulexusAssassinState, internal.EversorAssassinState, internal.VindicareAssassinState:
			object := record.SerializedObject.(*objects.AssassinState)
			units = append(units, object)
		}
	}

	sort.Sort(Sort(units))

	return units
}

func getClass(s string) string {
	class, _, _ := strings.Cut(s, "_")
	return class
}

func stem(s string) string {
	base, _, _ := strings.Cut(s, "_")
	return base
}

func mastercrafted(key string, trimPrefixes ...string) bool {
	for _, p := range trimPrefixes {
		key = strings.TrimPrefix(key, p)
	}
	return strings.Contains(key, "_")
}

var weaponRename = map[string]string{
	"Sword":  "ForceSword",
	"Shield": "StormShield",
	"Hammer": "DaemonHammer",
}

func getLvl(s string) int {
	_, after, _ := strings.Cut(s, "_")
	lvl, _ := strconv.Atoi(after)
	return lvl
}

type timelineUnlock struct {
	eventKey       string
	prereqID       string
	unlockedID     string
	requireKoramar bool
}

func (m *Manager) canUnlockTimelineEvent(u timelineUnlock) (enable, show bool) {
	var available, advancedTime bool
	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.GameUnlocksSaveState:
			object := record.SerializedObject.(*objects.GameUnlocksSaveState)
			for i := range object.Unlocks {
				switch object.Unlocks[i].ID {
				case u.prereqID:
					available = true
				case KoramarMissionDefeated:
					advancedTime = true
				case u.unlockedID:
					return false, true
				}
			}
		case internal.TimelineEventOccasion:
			object := record.SerializedObject.(*objects.TimelineEventOccasion)
			if object.EventToPlay.Key == u.eventKey {
				return object.TriggerTime > 0, object.TriggerTime == 0
			}
		}
	}

	if u.requireKoramar {
		return available && advancedTime, false
	}
	return available, false
}

func (m *Manager) unlockTimelineEvent(eventKey string, calendarType int, alsoReset ...string) {
	reset := func(o *objects.TimelineEventOccasion) {
		o.TriggerTime = 0
		o.SavedChosenResults.Values = []interface{}{}
	}

	var eventOccasion *objects.TimelineEventOccasion
	forEach(m, internal.TimelineEventOccasion, func(o *objects.TimelineEventOccasion) {
		if o.EventToPlay.Key == eventKey {
			eventOccasion = o
		}
	})

	if eventOccasion != nil {
		reset(eventOccasion)
		for _, key := range alsoReset {
			forEach(m, internal.TimelineEventOccasion, func(o *objects.TimelineEventOccasion) {
				if o.EventToPlay.Key == key {
					reset(o)
				}
			})
		}
		return
	}

	var saveState *objects.TimeManagerSaveState
	forEach(m, internal.TimeManagerSaveState, func(o *objects.TimeManagerSaveState) {
		saveState = o
	})
	if saveState == nil {
		return
	}

	id := m.generateNewInstanceId()

	eventOccasion = &objects.TimelineEventOccasion{}
	eventOccasion.EventToPlay.Key = eventKey
	eventOccasion.CalendarType = calendarType
	eventOccasion.SavedChosenResults.Values = []interface{}{}

	saveState.CurrentOccasions.Values = append(saveState.CurrentOccasions.Values, objects.IntValue{Key: id})
	m.state.LinearInstanceIds = append(m.state.LinearInstanceIds, id)
	m.state.LinearRecords = append(m.state.LinearRecords, &internal.LinearRecord{
		TypeName:         internal.TimelineEventOccasion,
		SerializedObject: eventOccasion,
	})
}

func (m *Manager) generateNewInstanceId() int {
	minId := 0
	for _, id := range m.state.LinearInstanceIds {
		if id < minId {
			minId = id
		}
	}
	return minId - 1
}

func forEach[T any](m *Manager, typeName string, fn func(*T)) {
	for _, record := range m.state.LinearRecords {
		if record.TypeName == typeName {
			fn(record.SerializedObject.(*T))
		}
	}
}

func first[T any](m *Manager, typeName string) *T {
	for _, record := range m.state.LinearRecords {
		if record.TypeName == typeName {
			return record.SerializedObject.(*T)
		}
	}
	return nil
}
