package main

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// OperationDefinition is what the generated fetch needs to know about the
// list operation, read from the Smithy model shipped with the SDK.
type OperationDefinition struct {
	Package         string
	FunctionName    string
	FunctionInput   string
	Item            string
	Items           string
	PageInputToken  string
	PageOutputToken string
	PageSize        string
}

func goGetPackage(pkg string) {
	cmd := exec.Command("go", "get", "github.com/aws/"+aws.SDKName+"/service/"+pkg)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

func readOperationDefinition(operation, pkg, functionName string) OperationDefinition {
	model := readSmithyModel(pkg)
	shape := jsonPath(model, "shapes", operation).(map[string]interface{})
	_, functionInput := parseOperation(jsonPath(shape, "input", "target").(string))

	definition := OperationDefinition{
		Package:       pkg,
		FunctionName:  functionName,
		FunctionInput: functionInput,
	}
	if paginated, isPaginated := jsonPath(shape, "traits", "smithy.api#paginated").(map[string]interface{}); isPaginated {
		definition.addPagination(model, paginated, jsonPath(shape, "output", "target").(string))
	}
	return definition
}

func (d *OperationDefinition) addPagination(model interface{}, paginated map[string]interface{}, outputShape string) {
	d.Items = paginated["items"].(string)
	itemsShape := jsonPath(model, "shapes", outputShape, "members", d.Items, "target").(string)
	_, d.Item = parseOperation(jsonPath(model, "shapes", itemsShape, "member", "target").(string))
	d.PageInputToken = optionalString(paginated["inputToken"])
	d.PageOutputToken = optionalString(paginated["outputToken"])
	d.PageSize = optionalString(paginated["pageSize"])
}

func optionalString(value interface{}) string {
	if value == nil {
		return ""
	}
	return value.(string)
}

// readSmithyModel loads the API model of the package from the SDK sources in
// the module cache.
func readSmithyModel(pkg string) interface{} {
	raw, err := os.ReadFile(smithyModelPath(pkg))
	if err != nil {
		panic(err)
	}
	var model interface{}
	if err := json.Unmarshal(raw, &model); err != nil {
		panic(err)
	}
	return model
}

func smithyModelPath(pkg string) string {
	glob := goPath() + "/pkg/mod/github.com/aws/" + aws.SDKName + "@v" + aws.SDKVersion + "/codegen/sdk-codegen/aws-models/" + pkg + ".*.json"
	matches, err := filepath.Glob(glob)
	if err != nil {
		panic(err)
	}
	switch len(matches) {
	case 0:
		panic(errors.New("Unable to find a file with glob \"" + glob + "\""))
	case 1:
		return matches[0]
	default:
		panic(errors.New("More than one file with glob \"" + glob + "\""))
	}
}

func goPath() string {
	if gopath, exists := os.LookupEnv("GOPATH"); exists {
		return gopath
	}
	userHome, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return userHome + "/go"
}

// jsonPath walks nested JSON objects along keys; it returns nil when a key
// is missing.
func jsonPath(document interface{}, keys ...string) interface{} {
	value := document
	for _, key := range keys {
		object, isObject := value.(map[string]interface{})
		if !isObject {
			return nil
		}
		value = object[key]
	}
	return value
}
