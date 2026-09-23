package searchers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type CognitoUserPoolSearcher struct{}

func (s CognitoUserPoolSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s CognitoUserPoolSearcher) fetch(cfg aws.Config) ([]types.UserPoolDescriptionType, error) {
	svc := cognitoidentityprovider.NewFromConfig(cfg)

	var entities []types.UserPoolDescriptionType
	nextToken := ""
	for {
		params := &cognitoidentityprovider.ListUserPoolsInput{
			MaxResults: aws.Int32(60), // max allowed by this API
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.ListUserPools(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.UserPools...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s CognitoUserPoolSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.UserPoolDescriptionType) {
	id := ""
	if entity.Id != nil {
		id = *entity.Id
	}
	title := id
	if entity.Name != nil && *entity.Name != "" {
		title = *entity.Name
	}

	subtitleArray := []string{id}
	if entity.LastModifiedDate != nil {
		subtitleArray = append(subtitleArray, "Modified "+entity.LastModifiedDate.Format(time.UnixDate))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	path := fmt.Sprintf("/cognito/v2/idp/user-pools/%s/users?region=%s", id, searchArgs.GetRegion())
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("cognito")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}
