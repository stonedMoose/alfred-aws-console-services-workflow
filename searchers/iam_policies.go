package searchers

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type IAMPolicySearcher struct{}

func (s IAMPolicySearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

// customer managed policies only; the AWS managed ones number in the thousands
// and are rarely what someone is looking for
func (s IAMPolicySearcher) fetch(cfg aws.Config) ([]types.Policy, error) {
	svc := iam.NewFromConfig(cfg)

	var entities []types.Policy
	marker := ""
	for {
		params := &iam.ListPoliciesInput{
			MaxItems: aws.Int32(1000), // get as many as we can
			Scope:    types.PolicyScopeTypeLocal,
		}
		if marker != "" {
			params.Marker = aws.String(marker)
		}
		resp, err := svc.ListPolicies(context.TODO(), params)
		if err != nil {
			return nil, err
		}

		entities = append(entities, resp.Policies...)

		if resp.IsTruncated && resp.Marker != nil {
			marker = *resp.Marker
		} else {
			break
		}
	}

	return entities, nil
}

func (s IAMPolicySearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.Policy) {
	title := *entity.PolicyName

	subtitleArray := []string{}
	subtitleArray = util.AppendString(subtitleArray, entity.Description)
	if entity.AttachmentCount != nil {
		subtitleArray = append(subtitleArray, "attached to "+strconv.Itoa(int(*entity.AttachmentCount)))
	}
	subtitleArray = util.AppendString(subtitleArray, entity.Path)
	subtitle := strings.Join(subtitleArray, " – ")

	arn := ""
	if entity.Arn != nil {
		arn = *entity.Arn
	}
	path := fmt.Sprintf("/iamv2/home#/policies/details/%s", url.QueryEscape(arn))
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("iam")).
		Valid(true)

	searchArgs.AddMatch(item, "arn:", arn, title)
}
