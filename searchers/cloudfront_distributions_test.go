package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestCloudFrontDistributionSearcher(t *testing.T) {
	TestSearcher(t, CloudFrontDistributionSearcher{}, util.GetCurrentFilename())
}
