package searchers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type StepFunctionsStateMachineSearcher struct{}

func (s StepFunctionsStateMachineSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "step_functions_state_machines", s.fetch, s.addToWorkflow)
}

func (StepFunctionsStateMachineSearcher) fetch(cfg aws.Config) ([]types.StateMachineListItem, error) {
	client := sfn.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.StateMachineListItem, *string, error) {
		resp, err := client.ListStateMachines(context.TODO(), &sfn.ListStateMachinesInput{
			MaxResults: 1000, // max allowed by this API
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.StateMachines, resp.NextToken, nil
	})
}

func (StepFunctionsStateMachineSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, stateMachine types.StateMachineListItem) {
	arn := aws.ToString(stateMachine.StateMachineArn)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       aws.ToString(stateMachine.Name),
		ConsolePath: "/states/home#/statemachines/view/" + arn,
		ServiceID:   "stepfunctions",
		ID:          arn,
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		string(stateMachine.Type),
		dateDetail("Created", stateMachine.CreationDate),
	))
}
