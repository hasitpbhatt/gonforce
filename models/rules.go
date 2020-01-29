package models

// Rules contains allowed and disallowed import path prefixes
type Rules struct {
	Allowed    []string `yaml:"allow"`
	Disallowed []string `yaml:"disallow"`
}
