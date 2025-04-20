package objects

type CurrencySaveState struct {
	SavedCurrencies []*Currency `json:"savedCurrencies"`
}

type Currency struct {
	CurrencyType StringValue `json:"currencyType"`
	Amount       int         `json:"amount"`
}
