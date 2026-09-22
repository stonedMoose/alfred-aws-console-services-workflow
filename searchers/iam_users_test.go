package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestIAMUserSearcher(t *testing.T) {
	TestSearcher(t, IAMUserSearcher{}, util.GetCurrentFilename())
}
