package regression_retrieval

import (
	"fmt"
	"os"
	"testing"

	"github.com/src-d/regression-retrieval/test"
	"github.com/src-d/regression-retrieval/test/gitcollector"
	metadataretrieval "github.com/src-d/regression-retrieval/test/metadata-retrieval"

	"github.com/src-d/regression-core"
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

// TestMetadataRetrieval
// runs regression comparison for remote:master and the first release
// <expected> no errors occurred during tests execution
// <expected> GetResults returns true
func TestMetadataRetrieval(t *testing.T) {
	config := regression.NewConfig()
	config.BinaryCache = "binaries"
	config.Versions = []string{"remote:master", "v0.1.0"}
	config.Repeat = 1
	// No token, no access to the v4 API
	config.GitHubToken = os.Getenv("REG_TOKEN")

	metadataRetrievalTest, err := test.NewTest(metadataretrieval.Kind, config)
	fmt.Printf("%+v", config)
	require.NoError(t, err)

	require.NoError(t, metadataRetrievalTest.Prepare())
	require.NoError(t, metadataRetrievalTest.RunLoad())

	metadataRetrievalTest.GetResults()
}
