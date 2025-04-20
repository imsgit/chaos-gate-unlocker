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

func (s Sort) Less(i, j int) bool {
	getClassStatusLvlAndName := func(obj interface{}) (class, status, lvl int, name string) {
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

	iClass, iStatus, iLvl, iName := getClassStatusLvlAndName(s[i])
	jClass, jStatus, jLvl, jName := getClassStatusLvlAndName(s[j])

	if jClass == iClass {
		if iStatus == jStatus {
			if jLvl == iLvl {
				return jName > iName
			}

			return jLvl < iLvl
		}

		return jStatus < iStatus
	}

	return jClass < iClass
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
	return strings.Split(s, "_")[0]
}

func getLvl(s string) int {
	lvl, _ := strconv.Atoi(strings.Split(s, "_")[1])
	return lvl
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
