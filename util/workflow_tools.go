package util

import (
	aw "github.com/deanishe/awgo"
)

// NewURLItem adds an item whose action opens a URL, with a cmd modifier that
// copies the URL instead.
func NewURLItem(wf *aw.Workflow, title string) *aw.Item {
	item := wf.NewItem(title).
		Valid(true).
		Var("action", "open-url")

	item.Cmd().Subtitle("Copy URL to clipboard").
		Var("action", "copy-to-clipboard")

	return item
}

// ServiceIcon returns the icon of an AWS service, shipped in the images folder
// under the service id.
func ServiceIcon(serviceID string) *aw.Icon {
	return &aw.Icon{Value: "images/" + serviceID + ".png"}
}
