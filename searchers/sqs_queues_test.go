package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestSQSQueueSearcher(t *testing.T) {
	TestSearcher(t, SQSQueueSearcher{}, util.GetCurrentFilename())
}
