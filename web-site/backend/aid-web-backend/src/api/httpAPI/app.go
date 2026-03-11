package httpAPI

var New app

type app struct{}

func (*app) Health() *HealthAPI { return &HealthAPI{} }
