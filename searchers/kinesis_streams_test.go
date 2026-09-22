package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestKinesisStreamSearcher(t *testing.T) {
	TestSearcher(t, KinesisStreamSearcher{}, util.GetCurrentFilename())
}
