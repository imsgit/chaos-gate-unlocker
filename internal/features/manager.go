package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"

	"cmp"
	"slices"
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

func (m *Manager) SetState(state *internal.State) {
	m.state = state
}

func classStatusLvlName(obj any) (class, status, lvl int, name string) {
	switch object := obj.(type) {
	case *objects.KnightState:
		status = object.HealthState.Status
		if object.CurrentSideMission.MissionID != "" {
			status = 5
		}
		return 2, status, getLvl(object.CurrentLevelData.Key), object.GivenName
	case *objects.AssassinState:
		return 1, 0, getLvl(object.CurrentLevelData.Key), object.GivenName
	case *objects.DreadnoughtState:
		return 0, 0, getLvl(object.CurrentLevelData.Key), object.GivenName
	default:
		return 0, 0, 0, ""
	}
}

func unitLess(a, b any) bool {
	iClass, iStatus, iLvl, iName := classStatusLvlName(a)
	jClass, jStatus, jLvl, jName := classStatusLvlName(b)
	return cmp.Or(
		cmp.Compare(jClass, iClass),
		cmp.Compare(jStatus, iStatus),
		cmp.Compare(jLvl, iLvl),
		strings.Compare(iName, jName),
	) < 0
}

func (m *Manager) Units() []any {
	var units []any
	knightInBarracks := map[int]bool{}

	for i, record := range m.state.LinearRecords {
		switch object := record.SerializedObject.(type) {
		case *objects.KnightsSaveState:
			for _, knight := range object.Knights {
				knightInBarracks[knight.Key] = true
			}
		case *objects.KnightState:
			if knightInBarracks[m.state.LinearInstanceIds[i]] {
				units = append(units, object)
			}
		case *objects.DreadnoughtState, *objects.AssassinState:
			units = append(units, object)
		}
	}

	sort.Slice(units, func(i, j int) bool { return unitLess(units[i], units[j]) })

	return units
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
		switch object := record.SerializedObject.(type) {
		case *objects.GameUnlocksSaveState:
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
		case *objects.TimelineEventOccasion:
			if object.EventToPlay.Key == u.eventKey {
				return object.TriggerTime > 0, object.TriggerTime == 0
			}
		}
	}

	return available && (advancedTime || !u.requireKoramar), false
}

func (m *Manager) unlockTimelineEvent(eventKey string, calendarType int, alsoReset ...string) {
	reset := func(o *objects.TimelineEventOccasion) {
		o.TriggerTime = 0
		o.SavedChosenResults.Values = []any{}
	}

	var eventOccasion *objects.TimelineEventOccasion
	forEach(m, func(o *objects.TimelineEventOccasion) {
		if o.EventToPlay.Key == eventKey {
			eventOccasion = o
		}
	})

	if eventOccasion != nil {
		reset(eventOccasion)
		forEach(m, func(o *objects.TimelineEventOccasion) {
			if slices.Contains(alsoReset, o.EventToPlay.Key) {
				reset(o)
			}
		})
		return
	}

	var saveState *objects.TimeManagerSaveState
	forEach(m, func(o *objects.TimeManagerSaveState) {
		saveState = o
	})
	if saveState == nil {
		return
	}

	id := m.generateNewInstanceId()

	eventOccasion = &objects.TimelineEventOccasion{}
	eventOccasion.EventToPlay.Key = eventKey
	eventOccasion.CalendarType = calendarType
	eventOccasion.SavedChosenResults.Values = []any{}

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
		minId = min(minId, id)
	}
	return minId - 1
}

func forEach[T any](m *Manager, fn func(*T)) {
	for _, record := range m.state.LinearRecords {
		if o, ok := record.SerializedObject.(*T); ok {
			fn(o)
		}
	}
}

func first[T any](m *Manager) *T {
	for _, record := range m.state.LinearRecords {
		if o, ok := record.SerializedObject.(*T); ok {
			return o
		}
	}
	return nil
}

func hasUnlock(o *objects.GameUnlocksSaveState, id string) bool {
	return slices.ContainsFunc(o.Unlocks, func(u objects.Unlock) bool { return u.ID == id })
}

func keysFromNames(dst []*objects.StringValue, names []string, byName map[string]string) []*objects.StringValue {
	dst = dst[:0]
	for _, name := range names {
		if name != "" {
			dst = append(dst, &objects.StringValue{Key: byName[name]})
		}
	}
	return dst
}

func lookupAt[V any](keys []*objects.StringValue, idx int, byID map[string]V) V {
	if idx >= 0 && idx < len(keys) {
		return byID[keys[idx].Key]
	}
	var zero V
	return zero
}
