package searchers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type CognitoIdentityPoolSearcher struct{}

func (s CognitoIdentityPoolSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s CognitoIdentityPoolSearcher) fetch(cfg aws.Config) ([]types.IdentityPoolShortDescription, error) {
	svc := cognitoidentity.NewFromConfig(cfg)

	var entities []types.IdentityPoolShortDescription
	nextToken := ""
	for {
		params := &cognitoidentity.ListIdentityPoolsInput{
			MaxResults: aws.Int32(60), // max allowed by this API
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.ListIdentityPools(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.IdentityPools...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s CognitoIdentityPoolSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.IdentityPoolShortDescription) {
	id := ""
	if entity.IdentityPoolId != nil {
		id = *entity.IdentityPoolId
	}
	title := id
	if entity.IdentityPoolName != nil && *entity.IdentityPoolName != "" {
		title = *entity.IdentityPoolName
	}

	path := fmt.Sprintf("/cognito/v2/identity/identity-pools/%s/user-statistics?region=%s", id, searchArgs.GetRegion())
	item := util.NewURLItem(wf, title).
		Subtitle(id).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("cognito")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}
