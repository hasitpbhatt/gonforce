package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"

	"github.com/hasitpbhatt/gonforce/models"
	yaml "gopkg.in/yaml.v2"
)

var _enforcerConfig enforcer

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
}

func main() {
	f, err := os.Open("gonforce.yaml")
	if err != nil {
		log.Fatal("File not found: gonforce.yaml")
	}
	defer f.Close()

	d := yaml.NewDecoder(f)
	d.SetStrict(true)
	if err := d.Decode(&_enforcerConfig); err != nil {
		log.Fatalf("Unable to decode gonforce.yaml: %v", err)
	}

	if err := _enforcerConfig.Default.Validate(); err != nil {
		log.Fatalf("Invalid gonforce.yaml: %v", err)
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Unable to get current dir: %v", err)
	}

	err = process(dir)
	if err != nil {
		log.Fatalf("The imports aren't as per the gonforce %v", err)
	}
}

func process(dir string) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ImportsOnly)
	if err != nil {
		log.Fatalf("Unable to parse imports: %v", err)
	}

	errorFound := false
	for _, pkg := range pkgs {
		for fpath, file := range pkg.Files {
			for _, imp := range file.Imports {
				if err := _enforcerConfig.Default.IsValidImport(fpath, imp.Path.Value); err != nil {
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
