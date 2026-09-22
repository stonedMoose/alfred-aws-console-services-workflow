package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestGlueJobSearcher(t *testing.T) {
	TestSearcher(t, GlueJobSearcher{}, util.GetCurrentFilename())
}
