package parsers

import (
	"log"
	"os"

	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"gopkg.in/yaml.v2"
)

// ParseConsoleServicesYml reads the catalogue of console services.
func ParseConsoleServicesYml(ymlPath string) []awsworkflow.AwsService {
	yamlFile, err := os.ReadFile(ymlPath)
	if err != nil {
		log.Fatal(err)
	}
	awsServices := []awsworkflow.AwsService{}
	if err := yaml.Unmarshal(yamlFile, &awsServices); err != nil {
		log.Fatal(err)
	}
	return awsServices
}
