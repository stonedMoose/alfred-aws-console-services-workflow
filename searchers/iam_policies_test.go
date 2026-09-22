package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestIAMPolicySearcher(t *testing.T) {
	TestSearcher(t, IAMPolicySearcher{}, util.GetCurrentFilename())
}
