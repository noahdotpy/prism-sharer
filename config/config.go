package config

type Config struct {
	InstancesDir string           `mapstructure:"instancesDir"`
	StoresDir    string           `mapstructure:"storesDir"`
	Groups       map[string]Group `mapstructure:"groups"`
}

type Group struct {
	Instances []string `mapstructure:"instances"`
	Resources []string `mapstructure:"resources"`
}
