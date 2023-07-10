package main_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	main "github.com/mjdusa/github-fork-update/cmd/github-fork-update"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Setup Suite
type GithubForkUpdateSuite struct {
	suite.Suite
}

func Test_GithubForkUpdate_Suite(t *testing.T) {
	suite.Run(t, &GithubForkUpdateSuite{})
}

func (s *GithubForkUpdateSuite) Test_GetUsage() {
	expected := fmt.Sprintf("usage:\n\t%s -auth='github-auth-token' [-verbose]\n\n", os.Args[0])
	actual := main.GetUsage()

	assert.Equal(s.T(), expected, actual, "GetUsage() message expected '%s', but got '%s'", expected, actual)
}

type TestGetParameters struct {
	Description     string
	TokenFlag       *string
	VerboseFlag     *bool
	ExpectedToken   string
	ExpectedVerbose bool
}

func (s *GithubForkUpdateSuite) Test_GetParameters() {
	expectedToken := "foo-bar"
	expectedVerboseFalse := false
	expectedVerboseTrue := true

	testList := []TestGetParameters{
		{
			Description:     "Default has no values",
			TokenFlag:       nil,
			VerboseFlag:     nil,
			ExpectedToken:   "",
			ExpectedVerbose: false,
		},
		{
			Description:     "Has only Token value",
			TokenFlag:       &expectedToken,
			VerboseFlag:     nil,
			ExpectedToken:   expectedToken,
			ExpectedVerbose: false,
		},
		{
			Description:     "Has only Verbose value false",
			TokenFlag:       nil,
			VerboseFlag:     &expectedVerboseFalse,
			ExpectedToken:   "",
			ExpectedVerbose: expectedVerboseFalse,
		},
		{
			Description:     "Has all values, Verbose value false",
			TokenFlag:       &expectedToken,
			VerboseFlag:     &expectedVerboseFalse,
			ExpectedToken:   expectedToken,
			ExpectedVerbose: expectedVerboseFalse,
		},
		{
			Description:     "Has only Verbose value true",
			TokenFlag:       nil,
			VerboseFlag:     &expectedVerboseTrue,
			ExpectedToken:   "",
			ExpectedVerbose: expectedVerboseTrue,
		},
		{
			Description:     "Has all values, Verbose value true",
			TokenFlag:       &expectedToken,
			VerboseFlag:     &expectedVerboseTrue,
			ExpectedToken:   expectedToken,
			ExpectedVerbose: expectedVerboseTrue,
		},
	}

	for _, test := range testList {
		os.Args = []string{"mainTest"}

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		if test.TokenFlag != nil {
			arg := fmt.Sprintf("-auth=%s", *test.TokenFlag)
			os.Args = append(os.Args, arg)
		}

		if test.VerboseFlag != nil && *test.VerboseFlag {
			os.Args = append(os.Args, "-verbose")
		}

		actualToken, actualVerbose := main.GetParameters()

		assert.Equal(s.T(), test.ExpectedToken, actualToken, "GetParameters() Token test '%s'", test.Description)
		assert.Equal(s.T(), test.ExpectedVerbose, actualVerbose, "GetParameters() Verbose test '%s'", test.Description)
	}
}
