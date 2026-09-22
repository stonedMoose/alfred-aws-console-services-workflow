package searchers

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type SecretsManagerSecretSearcher struct{}

func (s SecretsManagerSecretSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "secretsmanager_secrets", s.fetch, s.addToWorkflow)
}

func (SecretsManagerSecretSearcher) fetch(cfg aws.Config) ([]types.SecretListEntry, error) {
	client := secretsmanager.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.SecretListEntry, *string, error) {
		resp, err := client.ListSecrets(context.TODO(), &secretsmanager.ListSecretsInput{
			MaxResults: aws.Int32(100), // max allowed by this API
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.SecretList, resp.NextToken, nil
	})
}

func (SecretsManagerSecretSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, secret types.SecretListEntry) {
	name := aws.ToString(secret.Name)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title: name,
		// the console reads the region from the query string, next to the secret name
		ConsolePath: "/secretsmanager/secret?name=" + url.QueryEscape(name) + "&region=" + searchArgs.GetRegion(),
		ServiceID:   "secretsmanager",
		ID:          aws.ToString(secret.ARN),
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(
		aws.ToString(secret.Description),
		rotationDetail(secret.RotationEnabled),
		dateDetail("Changed", secret.LastChangedDate),
	))
}

func rotationDetail(rotationEnabled *bool) string {
	if aws.ToBool(rotationEnabled) {
		return "rotation enabled"
	}
	return ""
}
