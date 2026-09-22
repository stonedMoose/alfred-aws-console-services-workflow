package searchers

import (
	"testing"

	"github.com/rkoval/alfred-aws-console-services-workflow/util"
)

func TestStepFunctionsStateMachineSearcher(t *testing.T) {
	TestSearcher(t, StepFunctionsStateMachineSearcher{}, util.GetCurrentFilename())
}
