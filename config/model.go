package config

type Config struct {
	Name     string    `yaml:"name"`
	Commands []Command `yaml:"commands"`
}

type Command struct {
	Command string   `yaml:"cmd"`
	Syntax  string   `yaml:"run"`
	Args    []string `yaml:"args"`
}
