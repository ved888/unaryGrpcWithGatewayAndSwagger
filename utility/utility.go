package utility

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"log"
	"os"
)

func ConvertJSONToYML() {
	// Read the Swagger JSON file
	swaggerJSON, err := ioutil.ReadFile("doc/swagger/apidocs.swagger.json")
	if err != nil {
		log.Fatalf("Failed to read Swagger JSON file: %v", err)
	}
	// Parse the JSON
	var swaggerData interface{}
	err = json.Unmarshal(swaggerJSON, &swaggerData)
	if err != nil {
		log.Fatalf("Failed to parse Swagger JSON: %v", err)
	}
	// Convert to YAML
	swaggerYAML, err := yaml.Marshal(swaggerData)
	if err != nil {
		log.Fatalf("Failed to convert Swagger JSON to YAML: %v", err)
	}
	// Write the YAML to a file
	err = ioutil.WriteFile("doc/swagger/swagger.yaml", swaggerYAML, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to write Swagger YAML file: %v", err)
	}
	log.Println("Swagger JSON converted to YAML successfully!")
}
