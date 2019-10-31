package main

import (
	"os"

	"github.com/src-d/regression-retrieval/test"
	_ "github.com/src-d/regression-retrieval/test/gitcollector"
	_ "github.com/src-d/regression-retrieval/test/metadata-retrieval"

	"github.com/jessevdk/go-flags"
	"github.com/src-d/regression-core"
	"gopkg.in/src-d/go-log.v1"
)

var description = `data-retrieval utils regression tester.

This tool executes several versions of data-retrieval utils and compares query times and resource usage. There should be at least one version specified as an argument in the following way:

* v0.12.1 - release name from github (https://github.com/src-d/your-data-retrieval-util/releases). The binary will be downloaded.
* latest - latest release from github. The binary will be downloaded.
* remote:master - any tag or branch from repository. The binary will be built automatically.
* local:fix/some-bug - tag or branch from the repository in the current directory. The binary will be built.
* local:HEAD - current state of the repository. Binary is built.
* pull:266 - code from pull request #266 from your-data-retrieval-util repo. Binary is built.
* /path/to/your-data-retrieval-util - a binary built locally.

The repositories and downloaded/built binaries are cached by default in "repos" and "binaries" repositories from the current directory.
`

// Options CLI options
type Options struct {
	regression.Config

	CSV bool `long:"csv" description:"save csv files with last result"`

	Kind string `long:"kind" default:"gitcollector" description:"kind of utility that will be tested, currently only gitcollector is supported"`

	// prometheus pushgateway related options
	Prometheus bool `long:"prom" description:"store latest results to prometheus"`
	PromConfig regression.PromConfig
	CIConfig   regression.CIConfig
}

// TODO(@lwsanty) refactor to accept kind?
func main() {
	options := Options{
		Config: regression.NewConfig(),
	}

	parser := flags.NewParser(&options, flags.Default)
	parser.LongDescription = description

	args, err := parser.Parse()
	if err != nil {
		if err, ok := err.(*flags.Error); ok {
			if err.Type == flags.ErrHelp {
				os.Exit(0)
			}
		}

		log.Errorf(err, "Could not parse arguments")
		os.Exit(1)
	}

	config := options.Config

	if len(args) < 1 {
		log.Errorf(nil, "There should be at least one version")
		os.Exit(1)
	}

	config.Versions = args

	tst, err := test.NewTest(options.Kind, config)
	if err != nil {
		panic(err)
	}

	log.Infof("Preparing run")
	err = tst.Prepare()
	if err != nil {
		log.Errorf(err, "Could not prepare environment")
		os.Exit(1)
	}

	err = tst.RunLoad()
	if err != nil {
		panic(err)
	}

	tst.PrintTabbedResults()
	res := tst.GetResults()
	if !res {
		os.Exit(1)
	}
	if options.CSV {
		tst.SaveLatestCSV()
	}
	if options.Prometheus {
		if err := tst.StoreLatestToPrometheus(options.PromConfig, options.CIConfig); err != nil {
			log.Errorf(err, "Could not store results to prometheus")
			os.Exit(1)
		}
	}
}
