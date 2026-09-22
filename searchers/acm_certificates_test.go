package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestACMCertificateSearcher(t *testing.T) {
	TestSearcher(t, ACMCertificateSearcher{}, util.GetCurrentFilename())
}
