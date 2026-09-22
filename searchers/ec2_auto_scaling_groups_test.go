package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestEC2AutoScalingGroupSearcher(t *testing.T) {
	TestSearcher(t, EC2AutoScalingGroupSearcher{}, util.GetCurrentFilename())
}
