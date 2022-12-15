package config

type ConfigI interface {
	GetConfig() (*Config, error)
}
