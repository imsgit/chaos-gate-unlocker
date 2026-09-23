package features

import (
	"chaos-gate-unlocker/internal/objects"

	"slices"
)

const (
	Prognosticars        = "prognosticars"
	PrognosticarTutorial = "Prognosticar_Tutorial"
)

func prognosticarCurrency(o *objects.CurrencySaveState, current *objects.Currency) *objects.Currency {
	for _, c := range o.SavedCurrencies {
		if c.CurrencyType.Key == Prognosticars {
			current = c
		}
	}
	return current
}

func (m *Manager) RestorePrognosticars() {
	var currency *objects.Currency

	for _, record := range m.state.LinearRecords {
		switch object := record.SerializedObject.(type) {
		case *objects.GameUnlocksSaveState:
			object.Unlocks = slices.DeleteFunc(object.Unlocks, func(u objects.Unlock) bool { return u.ID == PrognosticarTutorial })
		case *objects.CurrencySaveState:
			currency = prognosticarCurrency(object, currency)
		case *objects.StarMapNodeModel:
			if object.HasPrognosticar.Value && currency != nil {
				object.HasPrognosticar.Value = false
				currency.Amount++
			}
		}
	}
}

func (m *Manager) CanRestorePrognosticars() (bool, bool) {
	var (
		currency  *objects.Currency
		available bool
		hasProg   bool
	)

	for _, record := range m.state.LinearRecords {
		switch object := record.SerializedObject.(type) {
		case *objects.GameUnlocksSaveState:
			if hasUnlock(object, PrognosticarTutorial) {
				available = true
			}
		case *objects.CurrencySaveState:
			currency = prognosticarCurrency(object, currency)
		case *objects.StarMapNodeModel:
			if object.HasPrognosticar.Value {
				hasProg = true
			}
		}
	}

	if hasProg {
		return true, false
	}

	if !available {
		return false, false
	}

	return false, currency != nil && currency.Amount > 0
}
