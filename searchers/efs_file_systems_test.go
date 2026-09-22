package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestEFSFileSystemSearcher(t *testing.T) {
	TestSearcher(t, EFSFileSystemSearcher{}, util.GetCurrentFilename())
}
