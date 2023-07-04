package config

type Config struct {
	InstancesDir string           `mapstructure:"instancesDir"`
	StoreDir     string           `mapstructure:"storeDir"`
	Groups       map[string]Group `mapstructure:"groups"`
}

type Group struct {
	Instances []string `mapstructure:"instances"`
	Resources []string `mapstructure:"resources"`
}
