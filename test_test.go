package regression_retrieval

import (
	"testing"

	"github.com/src-d/regression-core"
	"github.com/src-d/regression-retrieval/test"
	"github.com/src-d/regression-retrieval/test/gitcollector"
	"github.com/stretchr/testify/require"
)

// TestGitCollector
// runs regression comparison for remote:master remote:regression
// <expected> no errors occurred during tests execution
// <expected> GetResults returns true
func TestGitCollector(t *testing.T) {
	config := regression.NewConfig()
	config.BinaryCache = "binaries"
	config.Versions = []string{"remote:regression", "remote:master"}
	config.Repeat = 1

	gitCollectorTest, err := test.NewTest(gitcollector.Kind, config)
	require.NoError(t, err)

	require.NoError(t, gitCollectorTest.Prepare())
	require.NoError(t, gitCollectorTest.RunLoad())

	gitCollectorTest.GetResults()
}
