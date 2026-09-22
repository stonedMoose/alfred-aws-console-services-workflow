package searchers

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type DynamoDBTableSearcher struct{}

func (s DynamoDBTableSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "dynamodb_tables", s.fetch, s.addToWorkflow)
}

// fetch lists table names only, which is all ListTables returns. Describing
// every table for a richer subtitle would cost one API call per table, far too
// slow for a workflow that runs on each keystroke.
func (DynamoDBTableSearcher) fetch(cfg aws.Config) ([]string, error) {
	client := dynamodb.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(lastTableName string) ([]string, *string, error) {
		resp, err := client.ListTables(context.TODO(), &dynamodb.ListTablesInput{
			Limit:                   aws.Int32(100), // max allowed by this API
			ExclusiveStartTableName: awspaging.TokenOrNil(lastTableName),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.TableNames, resp.LastEvaluatedTableName, nil
	})
}

func (DynamoDBTableSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, tableName string) {
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       tableName,
		ConsolePath: "/dynamodbv2/home#table?name=" + url.QueryEscape(tableName),
		ServiceID:   "dynamodb",
	})
}
