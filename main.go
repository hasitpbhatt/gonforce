package main

import (
	"log"
	"os"

	"github.com/hasitpbhatt/gonforce/models"
	yaml "gopkg.in/yaml.v2"
)

type enforcer struct {
	//	Package is package name e.g. github.com/hasitpbhatt/gonforce
	Package string `yaml:"package"`

	//	Default rules contain the default allowed and not allowed packages
	//	e.g.
	//	default:
	//		allow:
	//			-	gopkg.in/yaml.v2
	//		disallow:
	//			-	gopkg.in/yaml.v1
	Default models.Rules `yaml:"default"`
}

func main() {
	f, err := os.Open("gonforce.yaml")
	if err != nil {
		log.Fatal("File not found: gonforce.yaml")
	}
	defer f.Close()

	e := enforcer{}
	d := yaml.NewDecoder(f)
	d.SetStrict(true)
	if err := d.Decode(&e); err != nil {
		log.Fatalf("Unable to decode gonforce.yaml: %v", err)
	}
}
