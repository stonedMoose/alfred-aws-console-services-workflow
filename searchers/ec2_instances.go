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

type EC2InstanceSearcher struct{}

func (s EC2InstanceSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "ec2_instances", s.fetch, s.addToWorkflow)
}

func (EC2InstanceSearcher) fetch(cfg aws.Config) ([]types.Instance, error) {
	client := ec2.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(pageToken string) ([]types.Instance, *string, error) {
		resp, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{
			MaxResults: aws.Int32(1000), // get as many as we can
			NextToken:  aws.String(pageToken),
		})
		if err != nil {
			return nil, nil, err
		}
		return instancesOf(resp.Reservations), resp.NextToken, nil
	})
}

func instancesOf(reservations []types.Reservation) []types.Instance {
	var instances []types.Instance
	for _, reservation := range reservations {
		instances = append(instances, reservation.Instances...)
	}
	return instances
}

func (EC2InstanceSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, instance types.Instance) {
	id := aws.ToString(instance.InstanceId)
	title, idDetail := namedOrID(ec2NameTag(instance.Tags), id)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       title,
		ConsolePath: "/ec2/v2/home#InstanceDetails:instanceId=" + id,
		ServiceID:   "ec2",
		ID:          id,
		IDPrefix:    "i-",
	}).Subtitle(joinKnown(" ", instanceStateEmoji(instance.State), idDetail, string(instance.InstanceType)))
}

var instanceStateEmojis = map[types.InstanceStateName]string{
	types.InstanceStateNameRunning:      "🟢",
	types.InstanceStateNameShuttingDown: "🟡",
	types.InstanceStateNameStopping:     "🟡",
	types.InstanceStateNameStopped:      "🔴",
	types.InstanceStateNameTerminated:   "🔴",
	types.InstanceStateNamePending:      "⚪️",
}

func instanceStateEmoji(state *types.InstanceState) string {
	if state == nil {
		return "❔"
	}
	if emoji, known := instanceStateEmojis[state.Name]; known {
		return emoji
	}
	return "❔"
}
