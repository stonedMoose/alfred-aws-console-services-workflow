package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestVPCSubnetSearcher(t *testing.T) {
	TestSearcher(t, VPCSubnetSearcher{}, util.GetCurrentFilename())
}
