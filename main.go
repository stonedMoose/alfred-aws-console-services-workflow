package main

import (
	"flag"
	"log"
	"strings"

	aw "github.com/deanishe/awgo"
	"github.com/deanishe/awgo/update"
	"github.com/rkoval/alfred-aws-console-services-workflow/workflow"
)

const (
	repository = "rkoval/alfred-aws-console-services-workflow"
	issuesURL  = "https://github.com/rkoval/alfred-aws-console-services-workflow/issues"
)

var (
	wf         *aw.Workflow
	forceFetch bool
	openAll    bool
	query      string
	ymlPath    string
)

func init() {
	flag.BoolVar(&forceFetch, "fetch", false, "force fetch via AWS instead of cache")
	flag.BoolVar(&openAll, "open_all", false, "open all URLs in a browser for the matching query")
	flag.StringVar(&query, "query", "", "query to use")
	flag.StringVar(&ymlPath, "yml_path", "console-services.yml", "path of the console services catalogue")
	flag.Parse()
	wf = aw.New(update.GitHub(repository), aw.HelpURL(issuesURL))
}

func main() {
	wf.Run(func() {
		log.Printf("running workflow with query: `%s`", query)
		workflow.Run(wf, strings.TrimLeft(query, " "), nil, forceFetch, openAll, ymlPath)
	})
}
