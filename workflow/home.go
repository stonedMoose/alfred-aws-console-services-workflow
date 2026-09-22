package workflow

import (
	"fmt"
	"log"

	aw "github.com/deanishe/awgo"
	"github.com/rkoval/alfred-aws-console-services-workflow/aliases"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsconfig"
	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

const (
	releasesURL = "https://github.com/rkoval/alfred-aws-console-services-workflow/releases"
	patreonURL  = "https://www.patreon.com/rkoval_alfred_aws_console_services_workflow"
	pizzaURL    = "https://ryankoval.pizza"
)

// showHome is what an empty query gets: a prompt, the profile and region in
// use, and a word from the author.
func (s *search) showHome() {
	log.Println("no search type parsed")
	s.wf.NewItem("Search for an AWS Service ...").
		Subtitle("e.g., cloudformation, ec2, s3 ...")
	s.addProfileStatusItem()
	s.addRegionStatusItem()
	addSupportItems(s.wf)
	checkForUpdate(s.wf)
	addUpdateAvailableItem(s.wf)
}

func (s *search) addProfileStatusItem() {
	if s.args.Profile == "" {
		addConfigurationHelpItem(s.wf, "No profile configured", awsconfig.ConfigFileDocsURL)
		return
	}
	s.wf.NewItem("Using profile \"" + s.args.Profile + "\"").
		Subtitle("Use \"" + aliases.OverrideAwsProfile + "\" to override for the current query").
		Icon(aw.IconAccount)
}

func (s *search) addRegionStatusItem() {
	if s.args.Cfg.Region == "" {
		addConfigurationHelpItem(s.wf, "No region configured for this profile", awsconfig.ConfigFileDocsURL)
		return
	}
	s.wf.NewItem("Using region \"" + s.args.Cfg.Region + "\"").
		Subtitle("Use \"" + aliases.OverrideAwsRegion + "\" to override for the current query").
		Icon(aw.IconWeb)
}

func addConfigurationHelpItem(wf *aw.Workflow, title, docsURL string) {
	util.NewURLItem(wf, title).
		Subtitle("Select this option to open AWS docs on how to configure").
		Arg(docsURL).
		Icon(aw.IconNote)
}

func addSupportItems(wf *aw.Workflow) {
	util.NewURLItem(wf, "Like this workflow? Consider donating! 😻").
		Subtitle("Select this option to open this project's Patreon").
		Arg(patreonURL).
		Icon(aw.IconFavorite)
	util.NewURLItem(wf, "Developers gotta eat! 🍕").
		Subtitle("Select this option to buy Ryan a pizza").
		Arg(pizzaURL).
		Icon(aw.IconColor)
}

func checkForUpdate(wf *aw.Workflow) {
	if !wf.UpdateCheckDue() {
		return
	}
	if err := wf.CheckForUpdate(); err != nil {
		wf.FatalError(err)
	}
}

func addUpdateAvailableItem(wf *aw.Workflow) {
	if !wf.UpdateAvailable() {
		return
	}
	util.NewURLItem(wf, fmt.Sprintf("Update available (current version: %s)", wf.Version())).
		Subtitle("Select this result to navigate to download").
		Arg(releasesURL).
		Icon(aw.IconInfo)
}
