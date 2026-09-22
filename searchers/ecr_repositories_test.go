package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestECRRepositorySearcher(t *testing.T) {
	TestSearcher(t, ECRRepositorySearcher{}, util.GetCurrentFilename())
}
