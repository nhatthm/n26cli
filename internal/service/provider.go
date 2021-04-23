package service

// ConfiguratorProvider provides configurator.Configurator.
type ConfiguratorProvider interface {
	Configurator() Configurator
}
