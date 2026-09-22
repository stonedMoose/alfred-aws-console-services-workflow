package searchers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type ElasticBeanstalkApplicationSearcher struct{}

func (s ElasticBeanstalkApplicationSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "elastic_beanstalk_applications", s.fetch, s.addToWorkflow)
}

func (ElasticBeanstalkApplicationSearcher) fetch(cfg aws.Config) ([]types.ApplicationDescription, error) {
	client := elasticbeanstalk.NewFromConfig(cfg)
	resp, err := client.DescribeApplications(context.TODO(), &elasticbeanstalk.DescribeApplicationsInput{})
	if err != nil {
		return nil, err
	}
	return resp.Applications, nil
}

func (ElasticBeanstalkApplicationSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, application types.ApplicationDescription) {
	arn := aws.ToString(application.ApplicationArn)
	name := aws.ToString(application.ApplicationName)
	title := arn
	if application.ApplicationName != nil {
		title = name
	}
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: "/elasticbeanstalk/home#/application/overview?applicationName=" + name,
		ServiceID:   "elasticbeanstalk",
		ID:          arn,
		IDPrefix:    searchutil.ARNPrefix,
	}).Subtitle(subtitleFrom(aws.ToString(application.Description)))
}
