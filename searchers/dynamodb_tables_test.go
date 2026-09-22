package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestDynamoDBTableSearcher(t *testing.T) {
	TestSearcher(t, DynamoDBTableSearcher{}, util.GetCurrentFilename())
}
