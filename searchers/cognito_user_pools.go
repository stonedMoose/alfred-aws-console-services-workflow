package searchers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type CognitoUserPoolSearcher struct{}

func (s CognitoUserPoolSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "cognito_user_pools", s.fetch, s.addToWorkflow)
}

func (CognitoUserPoolSearcher) fetch(cfg aws.Config) ([]types.UserPoolDescriptionType, error) {
	client := cognitoidentityprovider.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.UserPoolDescriptionType, *string, error) {
		resp, err := client.ListUserPools(context.TODO(), &cognitoidentityprovider.ListUserPoolsInput{
			MaxResults: aws.Int32(60), // max allowed by this API
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.UserPools, resp.NextToken, nil
	})
}

func (CognitoUserPoolSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, userPool types.UserPoolDescriptionType) {
	id := aws.ToString(userPool.Id)
	title, _ := namedOrID(aws.ToString(userPool.Name), id)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: fmt.Sprintf("/cognito/v2/idp/user-pools/%s/users?region=%s", id, searchArgs.GetRegion()),
		ServiceID:   "cognito",
	}).Subtitle(subtitleFrom(id, dateDetail("Modified", userPool.LastModifiedDate)))
}
