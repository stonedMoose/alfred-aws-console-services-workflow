package searchers

import (
	"context"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type EC2KeyPairSearcher struct{}

func (s EC2KeyPairSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "ec2_key_pairs", s.fetch, s.addToWorkflow)
}

// fetch makes a single call: DescribeKeyPairs returns every key pair at once
// and has no paging.
func (EC2KeyPairSearcher) fetch(cfg aws.Config) ([]types.KeyPairInfo, error) {
	client := ec2.NewFromConfig(cfg)
	resp, err := client.DescribeKeyPairs(context.TODO(), &ec2.DescribeKeyPairsInput{})
	if err != nil {
		return nil, err
	}
	return resp.KeyPairs, nil
}

func (EC2KeyPairSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, keyPair types.KeyPairInfo) {
	name := aws.ToString(keyPair.KeyName)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title: name,
		// key pairs have no detail page, so the list is opened filtered on the name
		ConsolePath: "/ec2/home#KeyPairs:search=" + url.QueryEscape(name),
		ServiceID:   "ec2",
	}).Subtitle(subtitleFrom(
		string(keyPair.KeyType),
		aws.ToString(keyPair.KeyPairId),
		aws.ToString(keyPair.KeyFingerprint),
	))
}
