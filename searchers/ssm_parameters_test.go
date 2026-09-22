package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestSSMParameterSearcher(t *testing.T) {
	TestSearcher(t, SSMParameterSearcher{}, util.GetCurrentFilename())
}
