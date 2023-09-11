package config

type ModulesMap map[string]*string

type Modules struct {
	Remote map[string]*string `yaml:"remote"`
	Local  map[string]*string `yaml:"local"`
}

type Config struct {
	Types   []string            `yaml:"types"`
	Bans    []map[string]string `yaml:"bans"`
	Modules Modules             `yaml:"modules"`
}
