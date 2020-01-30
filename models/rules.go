package models

// Rules contains allowed and disallowed import path prefixes
type Rules struct {
	Type    string   `yaml:"type"`
	Imports []string `yaml:"imports"`
	Except  []string `yaml:"except"`
}
