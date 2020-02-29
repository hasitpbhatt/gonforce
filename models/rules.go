package models

import (
	"fmt"
	"strings"
)

// RuleType signifies the type of the rule for given package
type RuleType string

const (
	// RuleTypeWhitelist showcases the imports mentioned
	// are whitelisted and all other (except "except" ones)
	// should not be used in the package
	RuleTypeWhitelist RuleType = "whitelist"

	// RuleTypeBlacklist showcases the imports mentioned
	// are blacklisted and all other (including the "except" ones)
	// can be used in the package
	RuleTypeBlacklist RuleType = "blacklist"
)

// PackageRule contains allowed and disallowed import path prefixes
type PackageRule struct {
	Type    RuleType `yaml:"type"`
	Imports []string `yaml:"imports"`
	Except  []string `yaml:"except"`
}

// Validate validates whether the Rules struct is valid or not
func (r PackageRule) Validate() error {
	if err := validateType(r.Type); err != nil {
		return err
	}
	if err := validateImports(r.Imports, r.Except); err != nil {
		return err
	}
	return nil
}

func validateType(t RuleType) error {
	switch t {
	case RuleTypeBlacklist, RuleTypeWhitelist:
		return nil
	default:
		return fmt.Errorf("Invalid rule type %v", t)
	}
}

func validateImports(imports, except []string) error {
	m := map[string]bool{}
	for _, path := range imports {
		m[path] = true
	}
	for _, path := range except {
		if m[path] {
			return fmt.Errorf("%s present in both imports and except", path)
		}
	}
	return nil
}

// IsValidImport checks whether the given import can be used in
// the given package or not
// It takes 2 arguments
// 1: file path that contains imports
// 2: import path being used
func (r PackageRule) IsValidImport(fpath, imp string) error {
	except := r.Except
	imports := r.Imports

	if r.Type == RuleTypeWhitelist {
		matched, isException := matches(except, imports, imp)
		if matched == "" {
			return fmt.Errorf("%v not whitelisted for %v", imp, fpath)
		}
		if isException {
			return fmt.Errorf("%v used in %v", imp, fpath)
		}
		return nil
	}

	matched, isAllowed := matches(except, imports, imp)
	if matched != "" && !isAllowed {
		return fmt.Errorf("%v not allowed for %v", imp, fpath)
	}
	return nil
}

func matches(set1, set2 []string, path string) (matched string, isSet1 bool) {
	for _, constraint := range set1 {
		if satisfies(path, constraint) {
			return constraint, true
		}
	}

	for _, constraint := range set2 {
		if satisfies(path, constraint) {
			return constraint, false
		}
	}

	return "", false
}

func satisfies(path, constraint string) bool {
	path = strings.Trim(path, "\"")
	if constraint == path {
		return true
	}
	return strings.HasPrefix(path+"/", constraint)
}
