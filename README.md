# regression-retrieval

**regression-retrieval** is a tool that runs different versions of data-retrieval utils and compares theirs resource consumption.

```
Usage:
  regression-retrieval [OPTIONS]

data-retrieval utils regression tester.

This tool executes several versions of data-retrieval utils and compares query times and resource usage. There should be at least two versions specified as arguments in the following way:

* v0.12.1 - release name from github (https://github.com/src-d/your-data-retrieval-util/releases). The binary will be downloaded.
* latest - latest release from github. The binary will be downloaded.
* remote:master - any tag or branch from repository. The binary will be built automatically.
* local:fix/some-bug - tag or branch from the repository in the current directory. The binary will be built.
* local:HEAD - current state of the repository. Binary is built.
* pull:266 - code from pull request #266 from your-data-retrieval-util repo. Binary is built.
* /path/to/your-data-retrieval-util - a binary built locally.

The repositories and downloaded/built binaries are cached by default in "repos" and "binaries" repositories from the current directory.


Application Options:
      --binaries=     Directory to store binaries (default: binaries) [$REG_BINARIES]
      --repos=        Directory to store repositories (default: repos) [$REG_REPOS]
      --url=          URL to the tool repo [$REG_GITURL]
      --gitport=      Port for local git server (default: 9418) [$REG_GITPORT]
      --repos-file=   YAML file with the list of repos [$REG_REPOS_FILE]
  -c, --complexity=   Complexity of the repositories to test (default: 1) [$REG_COMPLEXITY]
  -n, --repeat=       Number of times a test is run (default: 3) [$REG_REPEAT]
      --show-repos    List available repositories to test
  -t, --token=        Token used to connect to the API [$REG_TOKEN]
      --csv           save csv files with last result
      --kind=         kind of utility that will be tested, currently only gitcollector is supported
      --prom          store latest results to prometheus
      --prom-address= prometheus pushgateway address [$PROM_ADDRESS]
      --prom-job=     prometheus job [$PROM_JOB]
      --ci-branch=    branch env [$GIT_BRANCH]
      --ci-commit=    commit env [$GIT_COMMIT]

Help Options:
  -h, --help          Show this help message
```

#### Examples
##### Simple run
```
export GO111MODULE=on
export LOG_LEVEL=debug
export REG_TOKEN=your-github-token

cmd/regression-retrieval/regression-retrieval remote:master
```

##### Export to CSV and Prometheus pushgateway

```
export GO111MODULE=on
export PROM_ADDRESS="http://pushgateway-address:9091"
export PROM_JOB=retrieval_metrics
export GIT_BRANCH=master
export GIT_COMMIT=c2fc478d68af61e4f204623c4844e7c946b1ffd5
export REG_TOKEN=your-github-token
export LOG_LEVEL=debug

cmd/regression-retrieval/regression-retrieval --prom --csv remote:master
```