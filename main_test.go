package main

import (
	"flag"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestByGinkgo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Base Suite")
}

func parseTestCommandLineFlag() {
	flag.String("test.db.host", noDefaultValue, "postgreSQL hostname used in tests")
	if !flag.Parsed() {
		flag.Parse()
	}
}

func skipIfDatabaseIsNotSet(target func()) func() {
	parseTestCommandLineFlag()
	testDBHost := flag.Lookup("test.db.host").Value.String()
	testDBHost = strings.TrimRight(testDBHost, "/")
	if testDBHost == "" {
		skip := func() {
			Skip("Test flag is needed: -test.db.host=<host>")
		}
		return func() {
			BeforeEach(skip)
			target()
		}
	}
	testDBConnectionString = testDBHost + "/?sslmode=disable"
	testDBConnectionStringWithDatabase = testDBHost + "/test_hellofresh?sslmode=disable"
	return target
}
