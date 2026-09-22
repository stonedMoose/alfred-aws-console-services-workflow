package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestKMSKeySearcher(t *testing.T) {
	TestSearcher(t, KMSKeySearcher{}, util.GetCurrentFilename())
}
