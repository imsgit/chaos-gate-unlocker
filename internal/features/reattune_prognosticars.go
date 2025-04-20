package features

import (
	"chaos-gate-unlocker/internal"
	"chaos-gate-unlocker/internal/objects"
)

const (
	Prognosticars        = "prognosticars"
	PrognosticarTutorial = "Prognosticar_Tutorial"
)

func (m *Manager) ReattunePrognosticars() {
	var currency *objects.Currency

	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.GameUnlocksSaveState:
			object := record.SerializedObject.(*objects.GameUnlocksSaveState)
			temp := object.Unlocks[:0]
			for i := range object.Unlocks {
				if object.Unlocks[i].ID == PrognosticarTutorial {
					continue
				}
				temp = append(temp, object.Unlocks[i])
			}
			object.Unlocks = temp
		case internal.CurrencySaveState:
			object := record.SerializedObject.(*objects.CurrencySaveState)
			for i := range object.SavedCurrencies {
				if object.SavedCurrencies[i].CurrencyType.Key == Prognosticars {
					currency = object.SavedCurrencies[i]
				}
			}
		case internal.StarMapNodeModel:
			object := record.SerializedObject.(*objects.StarMapNodeModel)
			if object.HasPrognosticar.Value && currency != nil {
				object.HasPrognosticar.Value = false
				currency.Amount++
			}
		}
	}
}

func (m *Manager) CanReattunePrognosticars() (bool, bool) {
	var currency *objects.Currency

	for _, record := range m.state.LinearRecords {
		switch record.TypeName {
		case internal.CurrencySaveState:
			object := record.SerializedObject.(*objects.CurrencySaveState)
			for i := range object.SavedCurrencies {
				if object.SavedCurrencies[i].CurrencyType.Key == Prognosticars {
					currency = object.SavedCurrencies[i]
				}
			}
		case internal.StarMapNodeModel:
			object := record.SerializedObject.(*objects.StarMapNodeModel)
			if object.HasPrognosticar.Value {
				return true, false
			}
		}
	}

	return false, currency != nil && currency.Amount > 0
}
