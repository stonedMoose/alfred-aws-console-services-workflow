package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestCognitoIdentityPoolSearcher(t *testing.T) {
	TestSearcher(t, CognitoIdentityPoolSearcher{}, util.GetCurrentFilename())
}
