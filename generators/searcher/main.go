// Command searcher scaffolds a new searcher from the AWS API operation that
// lists its resources: the searcher file, its test, its registry entries and
// its workflow test cases.
//
//	go run ./generators/searcher Service Entity com.amazonaws.package#Operation
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/rkoval/alfred-aws-console-services-workflow/parsers"
)

func init() {
	configureAcronyms()
	flag.Parse()
}

// configureAcronyms teaches strcase the short names of the services, so that
// "SNS" stays "sns" in snake case rather than becoming "s_n_s".
func configureAcronyms() {
	for _, awsService := range parsers.ParseConsoleServicesYml("./console-services.yml") {
		if awsService.ShortName != "" {
			strcase.ConfigureAcronym(awsService.ShortName, strings.ToLower(awsService.ShortName))
		}
	}
}

func main() {
	args := flag.Args()
	if len(args) < 3 {
		usage()
	}
	service, entity, operation := args[0], args[1], args[2]

	pkg, functionName := parseOperation(operation)
	goGetPackage(pkg)
	definition := readOperationDefinition(operation, pkg, functionName)
	namer := NewSearcherNamer(service, entity, definition)

	appendToRegistry(namer)
	appendToWorkflowTest(namer)
	writeSearcherFile(namer)
	writeSearcherTestFile(namer)
}

func usage() {
	flag.Usage()
	fmt.Println("go run ./generators/searcher Service Entity com.amazonaws.package#FunctionName")
	os.Exit(1)
}

var operationNameRegex = regexp.MustCompile("com.amazonaws.([a-z0-9]+)#([a-zA-Z]+)")

// parseOperation splits "com.amazonaws.pkg#FunctionName" into its package and
// function name.
func parseOperation(operation string) (pkg, functionName string) {
	matches := operationNameRegex.FindStringSubmatch(operation)
	if len(matches) != 3 {
		log.Fatalln("operation argument must have the form \"com.amazonaws.pkg#FunctionName\"")
	}
	return matches[1], matches[2]
}

// SearcherNamer derives every name the generated files need from the service
// and entity names.
type SearcherNamer struct {
	ServiceLower      string
	EntityLowerPlural string
	StructName        string
	FileName          string
	OperationDefinition
}

var numberAfterUnderscore = regexp.MustCompile(`_([0-9]+)`)

// NewSearcherNamer expects entity in the singular, so that the plural forms
// can be derived from it.
func NewSearcherNamer(service, entity string, definition OperationDefinition) SearcherNamer {
	if strings.HasSuffix(entity, "s") {
		log.Fatalf("Entity should be singular for casing to work properly")
	}
	name := strings.ToTitle(service) + strings.ToTitle(entity)
	return SearcherNamer{
		ServiceLower:        strings.ToLower(service),
		EntityLowerPlural:   strings.ToLower(entity) + "s", // TODO make this proper english
		StructName:          name + "Searcher",
		FileName:            snakeCase(name) + "s", // TODO make this proper english
		OperationDefinition: definition,
	}
}

// snakeCase keeps digits attached to the word before them, as in "ec2",
// where strcase would start a new word.
func snakeCase(name string) string {
	return numberAfterUnderscore.ReplaceAllString(strcase.ToSnake(name), "$1")
}
