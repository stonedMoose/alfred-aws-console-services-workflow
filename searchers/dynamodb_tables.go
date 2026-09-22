package searchers

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type DynamoDBTableSearcher struct{}

func (s DynamoDBTableSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

// ListTables only returns table names, so that is all this searcher has to work
// with. Describing every table to build a richer subtitle would cost one API
// call per table, which is far too slow for a workflow that runs on each keystroke.
func (s DynamoDBTableSearcher) fetch(cfg aws.Config) ([]string, error) {
	svc := dynamodb.NewFromConfig(cfg)

	var entities []string
	exclusiveStartTableName := ""
	for {
		params := &dynamodb.ListTablesInput{
			Limit: aws.Int32(100), // max allowed by this API
		}
		if exclusiveStartTableName != "" {
			params.ExclusiveStartTableName = aws.String(exclusiveStartTableName)
		}
		resp, err := svc.ListTables(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.TableNames...)

		if resp.LastEvaluatedTableName != nil {
			exclusiveStartTableName = *resp.LastEvaluatedTableName
		} else {
			break
		}
	}

	return entities, nil
}

func (s DynamoDBTableSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity string) {
	title := entity

	path := fmt.Sprintf("/dynamodbv2/home#table?name=%s", url.QueryEscape(title))
	item := util.NewURLItem(wf, title).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("dynamodb")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}
