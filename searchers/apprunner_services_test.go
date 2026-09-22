package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestAppRunnerServiceSearcher(t *testing.T) {
	TestSearcher(t, AppRunnerServiceSearcher{}, util.GetCurrentFilename())
}
