package test

import (
	"github.com/src-d/regression-core"
	"gopkg.in/src-d/go-errors.v1"
)

// Constructor is a type that represents function of default Test Constructor
type Constructor func(config regression.Config) (Test, error)

var (
	// constructors is a map of all supported test constructors
	constructors = make(map[string]Constructor)

	errNotSupported = errors.NewKind("test kind %v is not supported")
)

// Test represents an interface for util-testing classes
type Test interface {
	// Prepare downloads, builds and runs entities for the  required environment
	Prepare() error
	// RunLoad does interaction with util/environment and obtains test results
	RunLoad() error
	// PrintTabbedResults prints obtained test results to stdout
	PrintTabbedResults()
	// SaveLatestCSV exports obtained test results to CSV files
	SaveLatestCSV()
	// StoreLatestToPrometheus pushes obtained test results to prometheus pushgateway
	StoreLatestToPrometheus(promConfig regression.PromConfig, ciConfig regression.CIConfig) error
	// GetResults compares two versions' test results and returns true if deviation is satisfying
	GetResults() bool
}

// Register updates the map of known test constructors
func Register(kind string, c Constructor) {
	constructors[kind] = c
}

// NewClient takes a given kind and creates related test
func NewTest(kind string, config regression.Config) (Test, error) {
	c, err := ValidateKind(kind)
	if err != nil {
		return nil, err
	}
	return c(config)
}

// ValidateKind checks if a given kind is supported
func ValidateKind(kind string) (Constructor, error) {
	c, ok := constructors[kind]
	if !ok {
		return nil, errNotSupported.New(kind)
	}

	return c, nil
}
