package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"strings"

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
	Default models.Rules `yaml:"default"`
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
		for fileString, file := range pkg.Files {
			for _, imp := range file.Imports {
				if err := isValid(fileString, imp.Path.Value); err != nil {
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

func isValid(fpath, imp string) error {
	except := _enforcerConfig.Default.Except
	imports := _enforcerConfig.Default.Imports

	if _enforcerConfig.Default.Type == "whitelist" {
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
	fmt.Println(path, constraint)
	path = strings.Trim(path, "\"")
	if constraint == path {
		return true
	}
	return strings.HasPrefix(path+"/", constraint)
}
