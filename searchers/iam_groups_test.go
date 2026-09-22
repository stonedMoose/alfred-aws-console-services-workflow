package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestIAMGroupSearcher(t *testing.T) {
	TestSearcher(t, IAMGroupSearcher{}, util.GetCurrentFilename())
}
