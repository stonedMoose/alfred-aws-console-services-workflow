package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestIAMRoleSearcher(t *testing.T) {
	TestSearcher(t, IAMRoleSearcher{}, util.GetCurrentFilename())
}
