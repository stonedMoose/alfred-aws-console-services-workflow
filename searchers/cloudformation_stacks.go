package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

type CloudFormationStackSearcher struct{}

func (s CloudFormationStackSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "cloudformation_stacks", s.fetch, s.addToWorkflow)
}

func (CloudFormationStackSearcher) fetch(cfg aws.Config) ([]types.Stack, error) {
	client := cloudformation.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.Stack, *string, error) {
		resp, err := client.DescribeStacks(context.TODO(), &cloudformation.DescribeStacksInput{
			NextToken: awspaging.TokenOrNil(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Stacks, resp.NextToken, nil
	})
}

// addToWorkflow does not go through AddConsoleResource: a stack matches on
// its Elastic Beanstalk environment name rather than on its generated name,
// and offers no autocompletion.
func (CloudFormationStackSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, stack types.Stack) {
	environmentName := cloudFormationTagValue(stack.Tags, "Name")
	title := stackTitle(aws.ToString(stack.StackName), environmentName)
	consolePath := "/cloudformation/home#/stacks/stackinfo?stackId=" + aws.ToString(stack.StackId)
	util.NewURLItem(wf, title).
		Subtitle(aws.ToString(stack.Description)).
		Arg(util.ConstructAWSConsoleUrl(consolePath, searchArgs.GetRegion())).
		// TODO: the stack icon should be "cloudformation"; "cloudwatch" is what the snapshots record
		Icon(util.ServiceIcon("cloudwatch")).
		Match(stackMatchTerm(title, environmentName, searchArgs.Query))
}

// stackTitle names the stacks Elastic Beanstalk generates after their environment.
func stackTitle(stackName, environmentName string) string {
	if isElasticBeanstalkStackName(stackName) && environmentName != "" {
		return fmt.Sprintf("%s (%s)", stackName, environmentName)
	}
	return stackName
}

// stackMatchTerm is what an Elastic Beanstalk stack is found by: its
// environment name, unless the user is typing the generated stack name.
func stackMatchTerm(title, environmentName, query string) string {
	if isElasticBeanstalkStackName(title) && !isElasticBeanstalkStackName(query) && environmentName != "" {
		return environmentName
	}
	return title
}

func isElasticBeanstalkStackName(name string) bool {
	return strings.HasPrefix(name, "awseb-")
}

func cloudFormationTagValue(tags []types.Tag, key string) string {
	for _, tag := range tags {
		if *tag.Key == key {
			return *tag.Value
		}
	}
	return ""
}
