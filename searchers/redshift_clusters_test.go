package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestRedshiftClusterSearcher(t *testing.T) {
	TestSearcher(t, RedshiftClusterSearcher{}, util.GetCurrentFilename())
}
