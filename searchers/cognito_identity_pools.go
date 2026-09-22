package searchers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type CognitoIdentityPoolSearcher struct{}

func (s CognitoIdentityPoolSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "cognito_identity_pools", s.fetch, s.addToWorkflow)
}

func (CognitoIdentityPoolSearcher) fetch(cfg aws.Config) ([]types.IdentityPoolShortDescription, error) {
	client := cognitoidentity.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.IdentityPoolShortDescription, *string, error) {
		resp, err := client.ListIdentityPools(context.TODO(), &cognitoidentity.ListIdentityPoolsInput{
			MaxResults: aws.Int32(60), // max allowed by this API
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.IdentityPools, resp.NextToken, nil
	})
}

func (CognitoIdentityPoolSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, identityPool types.IdentityPoolShortDescription) {
	id := aws.ToString(identityPool.IdentityPoolId)
	title, _ := namedOrID(aws.ToString(identityPool.IdentityPoolName), id)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: fmt.Sprintf("/cognito/v2/identity/identity-pools/%s/user-statistics?region=%s", id, searchArgs.GetRegion()),
		ServiceID:   "cognito",
	}).Subtitle(id)
}
