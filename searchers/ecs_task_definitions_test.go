package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestECSTaskDefinitionSearcher(t *testing.T) {
	TestSearcher(t, ECSTaskDefinitionSearcher{}, util.GetCurrentFilename())
}
