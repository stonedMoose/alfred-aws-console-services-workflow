package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type KMSKeySearcher struct{}

func (s KMSKeySearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

// aliases rather than ListKeys, because a bare key id is not something anyone
// searches for. AWS managed keys are dropped so this matches the console's
// customer managed keys page.
func (s KMSKeySearcher) fetch(cfg aws.Config) ([]types.AliasListEntry, error) {
	svc := kms.NewFromConfig(cfg)

	var entities []types.AliasListEntry
	marker := ""
	for {
		params := &kms.ListAliasesInput{
			Limit: aws.Int32(100), // max allowed by this API
		}
		if marker != "" {
			params.Marker = aws.String(marker)
		}
		resp, err := svc.ListAliases(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		for _, alias := range resp.Aliases {
			if alias.AliasName == nil || strings.HasPrefix(*alias.AliasName, "alias/aws/") {
				continue
			}
			if alias.TargetKeyId == nil {
				continue
			}
			entities = append(entities, alias)
		}

		if resp.Truncated && resp.NextMarker != nil {
			marker = *resp.NextMarker
		} else {
			break
		}
	}

	return entities, nil
}

func (s KMSKeySearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.AliasListEntry) {
	title := strings.TrimPrefix(*entity.AliasName, "alias/")

	subtitle := *entity.TargetKeyId

	path := fmt.Sprintf("/kms/home#/kms/keys/%s", *entity.TargetKeyId)
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("kms")).
		Valid(true)

	arn := ""
	if entity.AliasArn != nil {
		arn = *entity.AliasArn
	}
	searchArgs.AddMatch(item, "arn:", arn, title)
}
