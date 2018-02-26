package service

type HealthService struct{}

type HealthStatus struct {
	Alive   bool   `json:"alive"`
	Version string `json:"version"`
}

func (healthService *HealthService) Health() (HealthStatus, error) {
	health := HealthStatus{
		Alive:   true,
		Version: "0.0.1",
	}

	return health, nil
}
