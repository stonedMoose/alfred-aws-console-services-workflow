package searchers

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
	"github.com/rkoval/alfred-aws-console-services-workflow/caching"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type EC2KeyPairSearcher struct{}

func (s EC2KeyPairSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	cacheName := util.GetCurrentFilename()
	entities := caching.LoadEntityArrayFromCache(wf, searchArgs, cacheName, s.fetch)
	for _, entity := range entities {
		s.addToWorkflow(wf, searchArgs, entity)
	}
	return nil
}

// DescribeKeyPairs returns every key pair in one response — it has no paging
func (s EC2KeyPairSearcher) fetch(cfg aws.Config) ([]types.KeyPairInfo, error) {
	svc := ec2.NewFromConfig(cfg)

	resp, err := svc.DescribeKeyPairs(context.TODO(), &ec2.DescribeKeyPairsInput{})
	if err != nil {
		return nil, err
	}

	return resp.KeyPairs, nil
}

func (s EC2KeyPairSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, entity types.KeyPairInfo) {
	title := *entity.KeyName

	subtitleArray := []string{}
	if entity.KeyType != "" {
		subtitleArray = append(subtitleArray, string(entity.KeyType))
	}
	subtitleArray = util.AppendString(subtitleArray, entity.KeyPairId)
	subtitleArray = util.AppendString(subtitleArray, entity.KeyFingerprint)
	subtitle := strings.Join(subtitleArray, " – ")

	// key pairs have no detail page, so the list is opened filtered on the name
	path := fmt.Sprintf("/ec2/home#KeyPairs:search=%s", url.QueryEscape(title))
	item := util.NewURLItem(wf, title).
		Subtitle(subtitle).
		Arg(util.ConstructAWSConsoleUrl(path, searchArgs.GetRegion())).
		Icon(awsworkflow.GetImageIcon("ec2")).
		Valid(true)

	searchArgs.AddMatch(item, "", "", title)
}
