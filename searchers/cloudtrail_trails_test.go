package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestCloudTrailTrailSearcher(t *testing.T) {
	TestSearcher(t, CloudTrailTrailSearcher{}, util.GetCurrentFilename())
}
