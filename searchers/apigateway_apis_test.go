package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestAPIGatewayAPISearcher(t *testing.T) {
	TestSearcher(t, APIGatewayAPISearcher{}, util.GetCurrentFilename())
}
