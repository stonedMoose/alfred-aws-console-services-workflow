package searchers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type EC2ImageSearcher struct{}

func (s EC2ImageSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "ec2_amis", s.fetch, s.addToWorkflow)
}

// fetch is restricted to the images this account owns; the full public AMI
// catalog is hundreds of thousands of entries.
func (EC2ImageSearcher) fetch(cfg aws.Config) ([]types.Image, error) {
	client := ec2.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.Image, *string, error) {
		resp, err := client.DescribeImages(context.TODO(), &ec2.DescribeImagesInput{
			MaxResults: aws.Int32(1000), // get as many as we can
			Owners:     []string{"self"},
			NextToken:  awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Images, resp.NextToken, nil
	})
}

func (EC2ImageSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, image types.Image) {
	id := aws.ToString(image.ImageId)
	title, idDetail := namedOrID(aws.ToString(image.Name), id)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: "/ec2/home#ImageDetails:imageId=" + id,
		ServiceID:   "ec2",
	}).Subtitle(subtitleFrom(
		idDetail,
		aws.ToString(image.Description),
		string(image.Architecture),
		aws.ToString(image.CreationDate),
	))
}
