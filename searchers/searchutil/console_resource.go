package searchutil

import (
	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

// ARNPrefix is the IDPrefix of the resources that are identified by their ARN.
const ARNPrefix = "arn:"

// ConsoleResource is one AWS resource the way the workflow presents it: an
// Alfred item that opens the resource's own page in the console.
type ConsoleResource struct {
	// Title is what Alfred displays and matches on.
	Title string
	// ConsolePath is the path of the resource's console page; the region is
	// added when the URL is built.
	ConsolePath string
	// ServiceID names the owning service, which decides the icon.
	ServiceID string
	// ID is an identifier, typically the ARN, that may be typed in place of the title.
	ID string
	// IDPrefix reveals that the user is typing the identifier rather than the title.
	IDPrefix string
}

// AddConsoleResource adds the item of resource to the workflow and returns it
// for further decoration, such as a subtitle.
func (s *SearchArgs) AddConsoleResource(wf *aw.Workflow, resource ConsoleResource) *aw.Item {
	item := util.NewURLItem(wf, resource.Title).
		Arg(util.ConstructAWSConsoleUrl(resource.ConsolePath, s.GetRegion())).
		Icon(util.ServiceIcon(resource.ServiceID))
	return s.matchByIDOrTitle(item, resource.IDPrefix, resource.ID, resource.Title)
}
