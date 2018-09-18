package model

type ConfigGenerator interface {
}

// NewConfigGenerator creates a new instance of the dataplane configuration generator
func NewConfigGenerator(plugins []string) ConfigGenerator {
	return nil
}
