package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// DefaultPlaceholder used for setting sanitized fields.
const DefaultPlaceholder = "SANITIZED"

// Rules for configuring the dump.
type Rules struct {
	Sanitize map[string]Sanitize `yaml:"sanitize" json:"sanitize"`
	NoData   []string            `yaml:"nodata"   json:"nodata"`
	Ignore   []string            `yaml:"ignore"   json:"ignore"`
}

// Sanitize rules for while dumping a database.
type Sanitize map[string]string

// Load a config file.
func Load(path string) (Rules, error) {
	var rules Rules

	// We don't want to fail if the config file does not exist.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return rules, nil
	}

	f, err := ioutil.ReadFile(path)
	if err != nil {
		return rules, err
	}

	err = yaml.Unmarshal(f, &rules)
	if err != nil {
		return rules, err
	}

	return rules, nil
}
