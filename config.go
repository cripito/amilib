package amilib

type Config interface {
	Get() (*ConfigData, error)
	Put() error
	Delete() error
}

type ConfigData struct {
	NatsUrl  string
	UserName string
	Password string
	Server   string
	Port     int
	Verbose  bool
	Events   bool
}

func (c ConfigData) Get() (*ConfigData, error) {
	return &c, nil
}

func (c ConfigData) Put() error {
	return nil
}

func (c ConfigData) Delete() error {
	return nil
}

func NewConfig() *ConfigData {
	cfg := ConfigData{}

	return &cfg
}
