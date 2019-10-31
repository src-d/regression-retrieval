package metadataretrieval

import "github.com/src-d/regression-core"

// NewToolMetadataRetrieval creates a Tool with metadata-retrieval parameters filled
func NewToolMetadataRetrieval() regression.Tool {
	return regression.Tool{
		Name:        "metadata-retrieval",
		GitURL:      "https://github.com/src-d/metadata-retrieval",
		ProjectPath: "github.com/src-d/metadata-retrieval",
		BinaryName:  "cmd",
		BuildSteps: []regression.BuildStep{
			{
				Dir:     "",
				Command: "make",
				Args:    []string{"packages"},
				Env:     []string{"GOPROXY=https://proxy.golang.org"},
			},
		},
		// ExtraFiles: []string{
		// 	"testdata/regression.yml",
		// },
	}
}

// NewMetadataRetrieval returns a Binary struct for metadata-retrieval Tool
func NewMetadataRetrieval(
	config regression.Config,
	version string,
	releases *regression.Releases,
) *regression.Binary {
	return regression.NewBinary(config, NewToolMetadataRetrieval(), version, releases)
}
