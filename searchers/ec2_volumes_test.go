package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestEC2VolumeSearcher(t *testing.T) {
	TestSearcher(t, EC2VolumeSearcher{}, util.GetCurrentFilename())
}
