package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestCodeBuildProjectSearcher(t *testing.T) {
	TestSearcher(t, CodeBuildProjectSearcher{}, util.GetCurrentFilename())
}
