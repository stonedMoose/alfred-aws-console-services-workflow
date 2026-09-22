package searchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/awspaging"
	"github.com/rkoval/alfred-aws-console-services-workflow/searchers/searchutil"
)

// hostedZoneIDPrefix is what the SDK prepends to hosted zone ids and console
// links cannot have.
const hostedZoneIDPrefix = "/hostedzone/"

type Route53HostedZoneSearcher struct{}

func (s Route53HostedZoneSearcher) Search(wf *aw.Workflow, searchArgs searchutil.SearchArgs) error {
	return searchEntities(wf, searchArgs, "route53_hosted_zones", s.fetch, s.addToWorkflow)
}

func (Route53HostedZoneSearcher) fetch(cfg aws.Config) ([]types.HostedZone, error) {
	client := route53.NewFromConfig(cfg)
	return awspaging.FetchAllPages(func(marker string) ([]types.HostedZone, *string, error) {
		resp, err := client.ListHostedZones(context.TODO(), &route53.ListHostedZonesInput{
			MaxItems: aws.Int32(100),
			Marker:   awspaging.TokenOrNil(marker),
		})
		if err != nil {
			return nil, nil, err
		}
		return resp.HostedZones, resp.NextMarker, nil
	})
}

func (Route53HostedZoneSearcher) addToWorkflow(wf *aw.Workflow, searchArgs searchutil.SearchArgs, zone types.HostedZone) {
	id := aws.ToString(zone.Id)
	searchArgs.AddConsoleResource(wf, searchutil.ConsoleResource{
		Title:       aws.ToString(zone.Name),
		ConsolePath: "/route53/v2/hostedzones#ListRecordSets/" + strings.TrimPrefix(id, hostedZoneIDPrefix),
		ServiceID:   "route53",
		ID:          id,
		IDPrefix:    "Z",
	}).Subtitle(subtitleFrom(
		zoneVisibilityDetail(zone.Config),
		zoneCommentDetail(zone.Config),
		recordCountDetail(zone.ResourceRecordSetCount),
		id,
	))
}

func zoneVisibilityDetail(config *types.HostedZoneConfig) string {
	switch {
	case config == nil:
		return ""
	case config.PrivateZone:
		return "Private"
	default:
		return "Public"
	}
}

func zoneCommentDetail(config *types.HostedZoneConfig) string {
	if config == nil {
		return ""
	}
	return aws.ToString(config.Comment)
}

func recordCountDetail(count *int64) string {
	if count == nil {
		return ""
	}
	if *count > 1 {
		return fmt.Sprintf("%d records", *count)
	}
	return fmt.Sprintf("%d record", *count)
}
