package gitcollector

import "github.com/src-d/regression-core"

// NewToolGitCollector creates a Tool with gitcollector parameters filled
func NewToolGitCollector() regression.Tool {
	return regression.Tool{
		Name:        "gitcollector",
		GitURL:      "https://github.com/src-d/gitcollector",
		ProjectPath: "github.com/src-d/gitcollector",
		BuildSteps: []regression.BuildStep{
			{
				Dir:     "",
				Command: "make",
				Args:    []string{"packages"},
				Env:     []string{"GOPROXY=https://proxy.golang.org"},
			},
		},
		ExtraFiles: []string{
			"_testdata/regression.yml",
		},
	}
}

// NewGitCollector returns a Binary struct for gitcollector Tool
func NewGitCollector(
	config regression.Config,
	version string,
	releases *regression.Releases,
) *regression.Binary {
	return regression.NewBinary(config, NewToolGitCollector(), version, releases)
}
