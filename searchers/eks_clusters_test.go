package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestEKSClusterSearcher(t *testing.T) {
	TestSearcher(t, EKSClusterSearcher{}, util.GetCurrentFilename())
}
