package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestBatchJobQueueSearcher(t *testing.T) {
	TestSearcher(t, BatchJobQueueSearcher{}, util.GetCurrentFilename())
}
