package objects

type IntValue struct {
	Key int `json:"key"`
}

type StringValue struct {
	Key string `json:"key"`
}

type HealthState struct {
	Status                  int     `json:"status"`
	RecoveryTimeLeft        float64 `json:"recoveryTimeLeft"`
	HealingSuspended        bool    `json:"healingSuspended"`
	PartialRecoveryTimeLeft float64 `json:"partialRecoveryTimeLeft"`
}
