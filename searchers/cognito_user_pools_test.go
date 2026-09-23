package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestCognitoUserPoolSearcher(t *testing.T) {
	TestSearcher(t, CognitoUserPoolSearcher{}, util.GetCurrentFilename())
}
