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

type TestGetParameters struct {
	Description     string
	AuthFlag        *string
	DebugFlag       *bool
	VerboseFlag     *bool
	ExpectedAuth    string
	ExpectedDebug   bool
	ExpectedVerbose bool
}

func Call_GetParameters(s *GithubForkUpdateSuite) {
	os.Args = []string{"mainTest"}
	arg := fmt.Sprintf("-auth=%s", os.Getenv("GITHUB_AUTH"))
	os.Args = append(os.Args, arg)
	os.Args = append(os.Args, "-debug")
	os.Args = append(os.Args, "-verbose")

	actualAuth, actualDebug, actualVerbose := main.GetParameters()
	main.GetParameters()

	fmt.Println("inside")

	fmt.Fprintf(os.Stdout, "actualAuth=[%s]\n", actualAuth)
	fmt.Fprintf(os.Stdout, "actualDebug=[%t]\n", actualDebug)
	fmt.Fprintf(os.Stdout, "actualVerbose=[%t]\n", actualVerbose)
}

func (s *GithubForkUpdateSuite) Test_GetParameters() {
	ExpectedAuth := "foo-bar"
	expectedVerboseFalse := false
	expectedVerboseTrue := true

	testList := []TestGetParameters{
		{
			Description:     "Default has no values",
			AuthFlag:        nil,
			DebugFlag:       nil,
			VerboseFlag:     nil,
			ExpectedAuth:    "",
			ExpectedDebug:   false,
			ExpectedVerbose: false,
		},
		{
			Description:     "Has only Token value",
			AuthFlag:        &ExpectedAuth,
			DebugFlag:       nil,
			VerboseFlag:     nil,
			ExpectedAuth:    ExpectedAuth,
			ExpectedDebug:   false,
			ExpectedVerbose: false,
		},
		{
			Description:     "Has only Verbose value false",
			AuthFlag:        nil,
			DebugFlag:       nil,
			VerboseFlag:     &expectedVerboseFalse,
			ExpectedAuth:    "",
			ExpectedDebug:   false,
			ExpectedVerbose: expectedVerboseFalse,
		},
		{
			Description:     "Has all values, Verbose value false",
			AuthFlag:        &ExpectedAuth,
			DebugFlag:       nil,
			VerboseFlag:     &expectedVerboseFalse,
			ExpectedAuth:    ExpectedAuth,
			ExpectedDebug:   false,
			ExpectedVerbose: expectedVerboseFalse,
		},
		{
			Description:     "Has only Verbose value true",
			AuthFlag:        nil,
			DebugFlag:       nil,
			VerboseFlag:     &expectedVerboseTrue,
			ExpectedAuth:    "",
			ExpectedDebug:   false,
			ExpectedVerbose: expectedVerboseTrue,
		},
		{
			Description:     "Has all values, Verbose value true",
			AuthFlag:        &ExpectedAuth,
			DebugFlag:       nil,
			VerboseFlag:     &expectedVerboseTrue,
			ExpectedAuth:    ExpectedAuth,
			ExpectedDebug:   false,
			ExpectedVerbose: expectedVerboseTrue,
		},
	}

	for _, test := range testList {
		os.Args = []string{"mainTest"}

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		if test.AuthFlag != nil {
			arg := fmt.Sprintf("-auth=%s", *test.AuthFlag)
			os.Args = append(os.Args, arg)
		}

		if test.DebugFlag != nil && *test.DebugFlag {
			os.Args = append(os.Args, "-debug")
		}

		if test.VerboseFlag != nil && *test.VerboseFlag {
			os.Args = append(os.Args, "-verbose")
		}

		if test.AuthFlag == nil || len(*test.AuthFlag) == 0 || len(test.ExpectedAuth) == 0 {
			main.PanicOnExit = true

			defer func() {
				if r := recover(); r == nil {
					s.T().Errorf("The code did not panic")
				} else {
					s.T().Logf("Recovered in %v", r)
				}
			}()
		}

		actualAuth, actualDebug, actualVerbose := main.GetParameters()

		assert.Equal(s.T(), test.ExpectedAuth, actualAuth, "GetParameters() Auth test '%s'", test.Description)
		assert.Equal(s.T(), test.ExpectedDebug, actualDebug, "GetParameters() Debug test '%s'", test.Description)
		assert.Equal(s.T(), test.ExpectedVerbose, actualVerbose, "GetParameters() Verbose test '%s'", test.Description)
	}
}
