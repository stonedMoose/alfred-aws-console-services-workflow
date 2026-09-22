package searchers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk"
	"github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

type ElasticBeanstalkEnvironmentSearcher struct{}

func (s ElasticBeanstalkEnvironmentSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "elastic_beanstalk_environments", s.fetch, s.addToWorkflow)
}

func (ElasticBeanstalkEnvironmentSearcher) fetch(cfg aws.Config) ([]types.EnvironmentDescription, error) {
	client := elasticbeanstalk.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.EnvironmentDescription, *string, error) {
		resp, err := client.DescribeEnvironments(context.TODO(), &elasticbeanstalk.DescribeEnvironmentsInput{
			MaxRecords: aws.Int32(1000), // get as many as we can
			NextToken:  aws.String(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.Environments, resp.NextToken, nil
	})
}

func (ElasticBeanstalkEnvironmentSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, environment types.EnvironmentDescription) {
	id := aws.ToString(environment.EnvironmentId)
	applicationName := aws.ToString(environment.ApplicationName)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       aws.ToString(environment.EnvironmentName),
		ConsolePath: fmt.Sprintf("/elasticbeanstalk/home#/environment/%s?applicationName=%s&environmentId=%s", environmentConsolePage(environment), applicationName, id),
		ServiceID:   "elasticbeanstalk",
		ID:          id,
		IDPrefix:    "e-",
	}).Subtitle(joinKnown(" ", environmentHealthEmoji(environment.Health), id, applicationName))
}

// environmentConsolePage picks the events page for terminated environments,
// which have no dashboard any more.
func environmentConsolePage(environment types.EnvironmentDescription) string {
	if environment.Status == types.EnvironmentStatusTerminated {
		return "events"
	}
	return "dashboard"
}

var environmentHealthEmojis = map[types.EnvironmentHealth]string{
	types.EnvironmentHealthGreen:  "🟢",
	types.EnvironmentHealthYellow: "🟡",
	types.EnvironmentHealthRed:    "🔴",
	types.EnvironmentHealthGrey:   "⚪️",
}

func environmentHealthEmoji(health types.EnvironmentHealth) string {
	if emoji, known := environmentHealthEmojis[health]; known {
		return emoji
	}
	return "❔"
}
