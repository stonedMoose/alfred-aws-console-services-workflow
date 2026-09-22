package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestEC2ImageSearcher(t *testing.T) {
	TestSearcher(t, EC2ImageSearcher{}, util.GetCurrentFilename())
}
