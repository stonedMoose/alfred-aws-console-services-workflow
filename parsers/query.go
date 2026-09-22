package parsers

import (
	"strings"

	"github.com/rkoval/alfred-aws-console-services-workflow/awsconfig"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
)

// Query is what the user asked for: the service and sub-service named, the
// region or profile selected, and the words left over to search with.
type Query struct {
	RawQuery              string
	Service               *awsworkflow.AwsService
	SubService            *awsworkflow.AwsService
	HasTrailingWhitespace bool
	HasOpenAll            bool
	HasDefaultSearchAlias bool
	regionOverride        *awsconfig.Region
	// RegionQuery is set while the user is typing a region override.
	RegionQuery     *string
	ProfileOverride *awsconfig.Profile
	// ProfileQuery is set while the user is typing a profile override.
	ProfileQuery   *string
	RemainingQuery string
}

func (q *Query) IsEmpty() bool {
	return strings.Trim(q.RawQuery, " ") == ""
}

// GetRegionOverride returns the region selected in the query, which global
// services ignore.
func (q *Query) GetRegionOverride() *awsconfig.Region {
	if q.Service == nil || !q.Service.HasGlobalRegion {
		return q.regionOverride
	}
	return nil
}

// IsBareServiceId tells whether the query is a service id with nothing after
// it, like "ec2" as opposed to "ec2 ": the user may still be typing it.
func (q *Query) IsBareServiceId() bool {
	return q.Service != nil &&
		!q.HasTrailingWhitespace &&
		q.SubService == nil &&
		!q.HasDefaultSearchAlias &&
		q.RemainingQuery == ""
}

// TargetsResources tells whether the query asks for the resources of a
// service rather than for its list of sub-services: either through the search
// alias, or by naming a sub-service and going on typing.
func (q *Query) TargetsResources() bool {
	return q.HasDefaultSearchAlias ||
		(q.SubService != nil && (q.HasTrailingWhitespace || q.RemainingQuery != ""))
}
