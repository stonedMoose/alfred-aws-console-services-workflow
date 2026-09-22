package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestSecretsManagerSecretSearcher(t *testing.T) {
	TestSearcher(t, SecretsManagerSecretSearcher{}, util.GetCurrentFilename())
}
