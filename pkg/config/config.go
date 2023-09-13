package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Rules for configuring the dump.
type Rules struct {
	Rewrite map[string]Rewrite `yaml:"rewrite" json:"rewrite"`
	NoData  []string           `yaml:"nodata"  json:"nodata"`
	Ignore  []string           `yaml:"ignore"  json:"ignore"`
	Where   map[string]string  `yaml:"where"   json:"where"`
}

// Rewrite rules for while dumping a database.
type Rewrite map[string]string

// Load a config file.
func Load(path string) (Rules, error) {
	var rules Rules

	// We don't want to fail if the config file does not exist.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return rules, nil
	}

	f, err := os.ReadFile(path)
	if err != nil {
		return rules, err
	}

	err = yaml.Unmarshal(f, &rules)
	if err != nil {
		return rules, err
	}

	return rules, nil
}
