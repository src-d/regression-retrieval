package metadataretrieval

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/src-d/regression-retrieval/prometheus"
	"github.com/src-d/regression-retrieval/test"

	"github.com/src-d/regression-core"
	"gopkg.in/src-d/go-log.v1"
	"gopkg.in/yaml.v2"
)

type (
	metadataRetrievalResults map[string][]*Result
	versionResults           map[string]metadataRetrievalResults

	// Test holds the information about a metadata-retrieval test
	Test struct {
		config            regression.Config
		metadataRetrieval map[string]*regression.Binary
		// organizations is array of lists of coma-separated organizations
		organizations []string
		results       versionResults
		log           log.Logger
	}
)

// Kind is an identifier of util type to be tested, used in factory test constructor
const Kind = "metadata-retrieval"

func init() {
	test.Register(Kind, NewTest)
}

// Result is a wrapper around regression.Result that additionally contains organizations that were processed
type Result struct {
	*regression.Result
	Organizations string
}

// NewTest creates a new Test struct
func NewTest(config regression.Config) (test.Test, error) {
	l, err := (&log.LoggerFactory{Level: log.InfoLevel}).New(log.Fields{})
	if err != nil {
		return nil, err
	}

	return &Test{
		config: config,
		log:    l,
	}, nil
}

// Prepare downloads and builds required metadata-retrieval versions
func (t *Test) Prepare() error {
	return t.prepareMetadataRetrieval()
}

func (t *Test) prepareMetadataRetrieval() error {
	t.log.Infof("Preparing metadata-retrieval binaries")
	releases := regression.NewReleases("src-d", "metadata-retrieval", t.config.GitHubToken)

	t.metadataRetrieval = make(map[string]*regression.Binary, len(t.config.Versions))
	for _, version := range t.config.Versions {
		b := NewMetadataRetrieval(t.config, version, releases)
		err := b.Download()
		if err != nil {
			return err
		}

		t.metadataRetrieval[version] = b
	}

	return nil
}

// RunLoad executes the tests
func (t *Test) RunLoad() error {
	results := make(versionResults)

	for _, version := range t.config.Versions {
		_, ok := results[version]
		if !ok {
			results[version] = make(metadataRetrievalResults)
		}

		metadataRetrieval, ok := t.metadataRetrieval[version]
		if !ok {
			panic("metadataRetrieval not initialized. Was Prepare called?")
		}

		l := t.log.New(log.Fields{"version": version})

		l.Infof("Running version tests")

		times := t.config.Repeat
		if times < 1 {
			times = 1
		}

		// TODO(kyrcha): add a regression.yml file to metadata-retrieval
		t.organizations = []string{"git-fixtures"}

		for _, orgs := range t.organizations {
			results[version][orgs] = make([]*Result, times)
			for i := 0; i < times; i++ {
				l.New(log.Fields{
					"orgs": orgs,
				}).Infof("Running query")

				result, err := t.runLoadTest(metadataRetrieval, orgs)
				results[version][orgs][i] = result

				if err != nil {
					return err
				}
			}
		}
	}

	t.results = results

	return nil
}

// runLoadTest runs metadata-retrieval download command and saves execution time + memory usage
func (t *Test) runLoadTest(
	metadataRetrieval *regression.Binary,
	orgs string,
) (*Result, error) {
	t.log.Infof("Executing metadata-retrieval test")

	command := NewCommand(metadataRetrieval.Path, orgs)
	start := time.Now()
	if err := command.Run(map[string]string{
		"LOG_LEVEL":     "debug",
		"GITHUB_TOKENS": t.config.GitHubToken,
	}); err != nil {
		t.log.With(log.Fields{
			"orgs":               orgs,
			"metadata-retrieval": metadataRetrieval.Path,
		}).Errorf(err, "Could not execute metadata-retrieval")
		return nil, err
	}
	wall := time.Since(start)
	rusage := command.Rusage()

	t.log.With(log.Fields{
		"wall":   wall,
		"memory": rusage.Maxrss,
	}).Infof("Finished queries")

	result := &regression.Result{
		Wtime:  wall,
		Stime:  time.Duration(rusage.Stime.Nano()),
		Utime:  time.Duration(rusage.Utime.Nano()),
		Memory: rusage.Maxrss * 1024,
	}

	r := &Result{
		Result:        result,
		Organizations: orgs,
	}

	return r, nil
}

// PrintTabbedResults prints table with results to stdout
// Example:
// Org                 | remote:add-regression-config
// bblfsh,git-fixtures | 1m12.223594627s
func (t *Test) PrintTabbedResults() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.TabIndent|tabwriter.Debug)
	fmt.Fprint(w, "\x1b[1;33m Org \x1b[0m")
	versions := t.config.Versions
	for _, v := range versions {
		fmt.Fprintf(w, "\t\x1b[1;33m %s \x1b[0m", v)
	}
	fmt.Fprintf(w, "\n")

	for _, orgs := range t.organizations {
		fmt.Fprintf(w, "\x1b[1;37m %s \x1b[0m", orgs)
		var (
			mini    int
			min     time.Duration
			maxi    int
			max     time.Duration
			results []string
		)
		for i, v := range versions {
			if r, found := t.results[v][orgs]; !found {
				results = append(results, "--")
			} else {
				t := r[0].Wtime
				for _, ri := range r[1:] {
					if ri.Wtime < min {
						t = ri.Wtime
					}
				}

				if min == 0 {
					min = t
				}

				if max == 0 {
					max = t
				}

				if t < min {
					min = t
					mini = i
				}

				if t > max {
					max = t
					maxi = i
				}

				results = append(results, t.String())
			}
		}

		for i, r := range results {
			fmt.Fprint(w, "\t")
			if i == mini {
				fmt.Fprintf(w, "\x1b[1;32m %s \x1b[0m", r)
			} else if i == maxi {
				fmt.Fprintf(w, "\x1b[1;31m %s \x1b[0m", r)
			} else {
				fmt.Fprintf(w, "\x1b[1;37m %s \x1b[0m", r)
			}
		}
		fmt.Fprintf(w, "\n")
	}
	w.Flush()
	fmt.Println()
}

// GetResults prints test results and returns if the tests passed
func (t *Test) GetResults() bool {
	if len(t.config.Versions) < 1 {
		panic("there should be at least one version")
	}

	versions := t.config.Versions
	ok := true
	for i, version := range versions[0 : len(versions)-1] {
		fmt.Printf("%s - %s ####\n", version, versions[i+1])
		a := t.results[versions[i]]
		b := t.results[versions[i+1]]

		for _, orgs := range t.organizations {
			fmt.Printf("## Organizations {%s} ##\n", orgs)
			if _, found := a[orgs]; !found {
				fmt.Printf("# Skip - organizations {%s} not found for version: %s\n", orgs, versions[i])
				continue
			}
			if _, found := b[orgs]; !found {
				fmt.Printf("# Skip - organizations {%s} not found for version: %s\n", orgs, versions[i+1])
				continue
			}

			queryA := a[orgs][0]
			queryB := b[orgs][0]

			queryA.Result = average(a[orgs])
			queryB.Result = average(b[orgs])
			c := queryA.Result.ComparePrint(queryB.Result, 10.0)
			if !c {
				ok = false
			}
		}
	}

	return ok
}

// SaveLatestCSV saves test results in a CSV files
// created file examples:
//		- plot_org1_org2_org3_memory.csv
//		- plot_org1_org2_org3_time.csv
func (t *Test) SaveLatestCSV() {
	version := t.config.Versions[len(t.config.Versions)-1]
	for _, orgs := range t.organizations {
		res := average(t.results[version][orgs])
		if err := res.SaveAllCSV(fmt.Sprintf("plot_%s_", strings.Replace(orgs, ",", "_", -1))); err != nil {
			panic(err)
		}
	}
}

// StoreLatestToPrometheus stores latest version results to prometheus pushgateway
func (t *Test) StoreLatestToPrometheus(promConfig regression.PromConfig, ciConfig regression.CIConfig) error {
	version := t.config.Versions[len(t.config.Versions)-1]
	cli := prometheus.NewPromClient(Kind, promConfig)
	for _, orgs := range t.organizations {
		res := average(t.results[version][orgs])
		if err := cli.Dump(res, version, orgs, ciConfig.Branch, ciConfig.Commit); err != nil {
			return err
		}
	}
	return nil
}

func average(pr []*Result) *regression.Result {
	if len(pr) == 0 {
		return nil
	}

	results := make([]*regression.Result, 0, len(pr))
	for _, r := range pr {
		results = append(results, r.Result)
	}

	return regression.Average(results)
}

func loadOrganizationsYaml(file string) ([]string, error) {
	text, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var res []string
	err = yaml.Unmarshal(text, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
