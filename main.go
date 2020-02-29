package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/hasitpbhatt/gonforce/models"
	yaml "gopkg.in/yaml.v2"
)

var _enforcerConfig enforcer

// Rule contains the package name and rules that need to be applied
type Rule struct {
	Name        string             `yaml:"name"`
	PackageRule models.PackageRule `yaml:"rule"`
}

type enforcer struct {
	//	Package is package name e.g. github.com/hasitpbhatt/gonforce
	Package string `yaml:"package"`

	//	Default rules contain the default allowed and not allowed packages
	//	e.g.
	//	default:
	//		type: whitelist
	//		imports:
	//			-	gopkg.in
	//		except:
	//			-	gopkg.in/yaml.v2
	Default models.PackageRule `yaml:"default"`

	// Rules are array of rules containing package name and package
	// rule
	Rules []Rule `yaml:"rules"`
}

func main() {
	f, err := os.Open("gonforce.yaml")
	if err != nil {
		log.Fatal("File not found: gonforce.yaml")
	}
	defer f.Close()

	if err := decode(f); err != nil {
		log.Fatal(err)
	}

	if err := processRoot(); err != nil {
		log.Fatal(err)
	}
}

func decode(f io.Reader) error {
	d := yaml.NewDecoder(f)
	d.SetStrict(true)
	if err := d.Decode(&_enforcerConfig); err != nil {
		return fmt.Errorf("Unable to decode gonforce.yaml: %v", err)
	}

	if err := _enforcerConfig.Default.Validate(); err != nil {
		return fmt.Errorf("Invalid gonforce.yaml: %v", err)
	}

	for _, rule := range _enforcerConfig.Rules {
		if err := rule.PackageRule.Validate(); err != nil {
			return fmt.Errorf("Invalid gonforce.yaml: %v", err)
		}
	}
	return nil
}

func processRoot() error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Unable to get current dir: %v", err)
	}

	err = process(dir, _enforcerConfig.Default)
	if err != nil {
		return err
	}

	for _, rule := range _enforcerConfig.Rules {
		err := process(filepath.Join(dir, rule.Name), rule.PackageRule)
		if err != nil {
			return err
		}
	}
	return nil
}

func process(dir string, pr models.PackageRule) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ImportsOnly)
	if err != nil {
		return fmt.Errorf("Unable to parse imports in %s: %v", dir, err)
	}

	errorFound := false
	for _, pkg := range pkgs {
		for fpath, file := range pkg.Files {
			for _, imp := range file.Imports {
				if err := pr.IsValidImport(fpath, imp.Path.Value); err != nil {
					errorFound = true
					fmt.Println(err)
				}
			}
		}
	}
	if errorFound {
		return fmt.Errorf("validation failed in %v", dir)
	}
	return nil
}
