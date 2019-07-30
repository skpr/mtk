package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// DefaultPlaceholder used for setting sanitized fields.
const DefaultPlaceholder = "SANITIZED"

// File for configuring dump rules.
type File struct {
	Sanitize Sanitize `yaml:"sanitize" json:"sanitize"`
	NoData   []string `yaml:"nodata"   json:"nodata"`
	Ignore   []string `yaml:"ignore"   json:"ignore"`
}

// Sanitize rules for while dumping a database.
type Sanitize struct {
	Tables []Table `yaml:"tables" json:"tables"`
}

// Table rules for while dumping a database.
type Table struct {
	Name   string  `yaml:"name"   json:"name"`
	Fields []Field `yaml:"fields" json:"fields"`
}

// Field rules for while dumping a database.
type Field struct {
	Name  string `yaml:"name"  json:"name"`
	Value string `yaml:"value" json:"value"`
}

// Load a config file.
func Load(path string) (File, error) {
	var file File

	// We don't want to fail if the config file does not exist.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return file, nil
	}

	f, err := ioutil.ReadFile(path)
	if err != nil {
		return file, err
	}

	err = yaml.Unmarshal(f, &file)
	if err != nil {
		return file, err
	}

	return file, nil
}
