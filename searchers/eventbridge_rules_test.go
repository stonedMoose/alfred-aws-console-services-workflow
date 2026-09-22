package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestEventBridgeRuleSearcher(t *testing.T) {
	TestSearcher(t, EventBridgeRuleSearcher{}, util.GetCurrentFilename())
}
