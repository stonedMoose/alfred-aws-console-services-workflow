package searchers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type StepFunctionsStateMachineSearcher struct{}

func (s StepFunctionsStateMachineSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s StepFunctionsStateMachineSearcher) fetch(cfg aws.Config) ([]types.StateMachineListItem, error) {
	svc := sfn.NewFromConfig(cfg)

	var entities []types.StateMachineListItem
	nextToken := ""
	for {
		params := &sfn.ListStateMachinesInput{
			MaxResults: 1000, // max allowed by this API
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.ListStateMachines(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.StateMachines...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s StepFunctionsStateMachineSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.StateMachineListItem) {
	title := *entity.Name

	subtitleArray := []string{}
	if entity.Type != "" {
		subtitleArray = append(subtitleArray, string(entity.Type))
	}
	if entity.CreationDate != nil {
		subtitleArray = append(subtitleArray, "Created "+entity.CreationDate.Format(time.UnixDate))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/states/home#/statemachines/view/%s", *entity.StateMachineArn)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("stepfunctions")).
		Valid(true)

	searchArgs.AddMatch(item, "arn:", *entity.StateMachineArn, title)
}
