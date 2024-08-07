package amilib

type Config interface {
	Get() (*ConfigData, error)
	Put() error
}

type ConfigData struct {
}
