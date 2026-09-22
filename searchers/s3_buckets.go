package searchers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type S3BucketSearcher struct{}

func (s S3BucketSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "s3_buckets", s.fetch, s.addToWorkflow)
}

func (S3BucketSearcher) fetch(cfg aws.Config) ([]types.Bucket, error) {
	client := s3.NewFromConfig(cfg)
	resp, err := client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}
	return resp.Buckets, nil
}

func (S3BucketSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, bucket types.Bucket) {
	name := aws.ToString(bucket.Name)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title: name,
		// the region goes in the path because the console reads it from the
		// query string next to the bucket name
		ConsolePath: fmt.Sprintf("/s3/buckets/%s/?region=%s&tab=objects", name, searchArgs.Cfg.Region),
		ServiceID:   "s3",
	}).Subtitle(dateDetail("Created", bucket.CreationDate))
}
