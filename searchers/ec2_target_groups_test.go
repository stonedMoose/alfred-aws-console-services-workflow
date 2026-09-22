package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestEC2TargetGroupSearcher(t *testing.T) {
	TestSearcher(t, EC2TargetGroupSearcher{}, util.GetCurrentFilename())
}
