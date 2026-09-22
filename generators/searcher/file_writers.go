package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"
)

const (
	registryFile     = "searchers/registry.go"
	workflowTestFile = "workflow/workflow_test.go"
)

// lastMapEntryRegex finds the end of the last entry of a map literal, where a
// new entry can be inserted.
var lastMapEntryRegex = regexp.MustCompile(`(,\n)(})`)

// appendToRegistry registers the searcher under its sub-service id and, when
// the service has no searcher yet, under the bare service id as well.
func appendToRegistry(namer SearcherNamer) {
	subServiceKey := namer.ServiceLower + "_" + namer.EntityLowerPlural
	content := addRegistryEntry(subServiceKey, namer.StructName)
	if content == "" {
		return
	}
	if !strings.Contains(content, fmt.Sprintf("\"%s\"", namer.ServiceLower)) {
		addRegistryEntry(namer.ServiceLower, namer.StructName)
	}
}

func addRegistryEntry(key, structName string) string {
	replacement := fmt.Sprintf("$1\t\"%s\": %s{}$1$2", key, structName)
	return modifyFileWithRegexReplace(registryFile, lastMapEntryRegex, replacement, "\""+key+"\"")
}

// appendToWorkflowTest adds the queries reaching the new searcher to the
// workflow test, reusing the searcher's own fixture.
func appendToWorkflowTest(namer SearcherNamer) {
	lastTestCaseRegex := regexp.MustCompile(`(\t},)(\n})`)
	addTestCase := func(query string) {
		testCase := fmt.Sprintf(`	{
		query:       "%s",
		fixtureName: "../searchers/%s_test", // reuse test fixture from this other test
	},`, query, namer.FileName)
		modifyFileWithRegexReplace(workflowTestFile, lastTestCaseRegex, "$1\n"+testCase+"$2", "\""+query+"\"")
	}

	service := namer.ServiceLower
	addTestCase(service)
	addTestCase(service + " ")
	addTestCase(service + " " + namer.EntityLowerPlural)
	addTestCase(service + " " + namer.EntityLowerPlural + " ")
}

const searcherTemplate = `package searchers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/{{ .Package }}"
	"github.com/aws/aws-sdk-go-v2/service/{{ .Package }}/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type {{ .StructName }} struct{}

func (s {{ .StructName }}) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "{{ .FileName }}", s.fetch, s.addToWorkflow)
}

func ({{ .StructName }}) fetch(cfg aws.Config) ([]types.{{ .Item }}, error) {
	client := {{ .Package }}.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.{{ .Item }}, *string, error) {
		resp, err := client.{{ .FunctionName }}(context.TODO(), &{{ .Package }}.{{ .FunctionInput }}{
			{{- if .PageSize }}
			{{ .PageSize }}: aws.Int32(100), // TODO max allowed by this API
			{{- end }}
			{{- if .PageInputToken }}
			{{ .PageInputToken }}: awspaging.TokenOrNil(pageToken),
			{{- end }}
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.{{ .Items }}, {{ if .PageOutputToken }}resp.{{ .PageOutputToken }}{{ else }}nil{{ end }}, nil
	})
}

func ({{ .StructName }}) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.{{ .Item }}) {
	title := aws.ToString(entity.TODO)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: "/{{ .ServiceLower }}/{{ .EntityLowerPlural }}/" + title, // TODO check the console path
		ServiceID:   "{{ .ServiceLower }}",
		ID:          aws.ToString(entity.TODOArn),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(aws.ToString(entity.TODO)))
}
`

const searcherTestTemplate = `package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func Test{{ .StructName }}(t *testing.T) {
	TestSearcher(t, {{ .StructName }}{}, util.GetCurrentFilename())
}
`

func writeSearcherFile(namer SearcherNamer) {
	writeTemplateToFile("searcher_file", searcherTemplate, fmt.Sprintf("searchers/%s.go", namer.FileName), namer)
}

func writeSearcherTestFile(namer SearcherNamer) {
	writeTemplateToFile("searcher_test_file", searcherTestTemplate, fmt.Sprintf("searchers/%s_test.go", namer.FileName), namer)
}

func writeTemplateToFile(templateName, templateString, fileName string, data interface{}) {
	t, err := template.New(templateName).Parse(templateString)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := t.Execute(file, data); err != nil {
		panic(err)
	}
}

// modifyFileWithRegexReplace rewrites the file with the regex replaced, and
// returns the new content; it leaves the file alone and returns "" when the
// file already contains ignoreIfContains.
func modifyFileWithRegexReplace(filename string, regex *regexp.Regexp, replacement string, ignoreIfContains string) string {
	raw, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	content := string(raw)
	if ignoreIfContains != "" && strings.Contains(content, ignoreIfContains) {
		return ""
	}
	replaced := regex.ReplaceAllString(content, replacement)
	if err := os.WriteFile(filename, []byte(replaced), 0600); err != nil {
		panic(err)
	}
	return replaced
}
