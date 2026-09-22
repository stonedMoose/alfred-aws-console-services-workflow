package searchers

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type SecretsManagerSecretSearcher struct{}

func (s SecretsManagerSecretSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

func (s SecretsManagerSecretSearcher) fetch(cfg aws.Config) ([]types.SecretListEntry, error) {
	svc := secretsmanager.NewFromConfig(cfg)

	var entities []types.SecretListEntry
	nextToken := ""
	for {
		params := &secretsmanager.ListSecretsInput{
			MaxResults: aws.Int32(100), // max allowed by this API
		}
		if nextToken != "" {
			params.NextToken = aws.String(nextToken)
		}
		resp, err := svc.ListSecrets(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.SecretList...)

		if resp.NextToken != nil && *resp.NextToken != "" {
			nextToken = *resp.NextToken
		} else {
			break
		}
	}

	return entities, nil
}

func (s SecretsManagerSecretSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.SecretListEntry) {
	title := *entity.Name

	subtitleArray := []string{}
	subtitleArray = util.AppendString(subtitleArray, entity.Description)
	if entity.RotationEnabled != nil && *entity.RotationEnabled {
		subtitleArray = append(subtitleArray, "rotation enabled")
	}
	if entity.LastChangedDate != nil {
		subtitleArray = append(subtitleArray, "Changed "+entity.LastChangedDate.Format(time.UnixDate))
	}
	subtitle := strings.Join(subtitleArray, " – ")

	// the region is written into the path here because the console reads it from
	// the query string alongside the secret name
	path := fmt.Sprintf("/secretsmanager/secret?name=%s&region=%s", url.QueryEscape(title), searchArgs.GetRegion())
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("secretsmanager")).
		Valid(true)

	arn := ""
	if entity.ARN != nil {
		arn = *entity.ARN
	}
	searchArgs.AddMatch(item, "arn:", arn, title)
}
