package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestEventBridgeEventBusSearcher(t *testing.T) {
	TestSearcher(t, EventBridgeEventBusSearcher{}, util.GetCurrentFilename())
}
