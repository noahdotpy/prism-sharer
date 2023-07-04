package config

type Config struct {
	InstancesDir string           `mapstructure:"instancesDir"`
	Groups       map[string]Group `mapstructure:"groups"`
}
type Group struct {
	Instances []string `mapstructure:"instances"`
	Shared    []string `mapstructure:"shared"`
}
