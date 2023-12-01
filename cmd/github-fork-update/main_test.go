package main

import (
	"flag"
	"fmt"
	"os"
	"testing"

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
	DebugFlag       bool
	VerboseFlag     bool
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

	actualAuth, actualDebug, actualVerbose := GetParameters()
	GetParameters()

	fmt.Println("inside")

	fmt.Fprintf(os.Stdout, "actualAuth=[%s]\n", actualAuth)
	fmt.Fprintf(os.Stdout, "actualDebug=[%t]\n", actualDebug)
	fmt.Fprintf(os.Stdout, "actualVerbose=[%t]\n", actualVerbose)
}

func (s *GithubForkUpdateSuite) Test_GetParameters() {
	ExpectedAuth := "foo-bar"
	expectedFalse := false
	expectedTrue := true

	testList := []TestGetParameters{
		{
			Description:     "Has only Token value",
			AuthFlag:        &ExpectedAuth,
			DebugFlag:       expectedFalse,
			VerboseFlag:     expectedFalse,
			ExpectedAuth:    ExpectedAuth,
			ExpectedDebug:   expectedFalse,
			ExpectedVerbose: expectedFalse,
		},
		{
			Description:     "Has all values, Debug value false",
			AuthFlag:        &ExpectedAuth,
			DebugFlag:       expectedFalse,
			VerboseFlag:     expectedFalse,
			ExpectedAuth:    ExpectedAuth,
			ExpectedDebug:   expectedFalse,
			ExpectedVerbose: expectedFalse,
		},
		{
			Description:     "Has all values, Verbose value true",
			AuthFlag:        &ExpectedAuth,
			DebugFlag:       expectedFalse,
			VerboseFlag:     expectedTrue,
			ExpectedAuth:    ExpectedAuth,
			ExpectedDebug:   expectedFalse,
			ExpectedVerbose: expectedTrue,
		},
		{
			Description:     "Has all values, Debug and Verbose value true",
			AuthFlag:        &ExpectedAuth,
			DebugFlag:       expectedTrue,
			VerboseFlag:     expectedTrue,
			ExpectedAuth:    ExpectedAuth,
			ExpectedDebug:   expectedTrue,
			ExpectedVerbose: expectedTrue,
		},
	}

	for _, test := range testList {
		os.Args = []string{"mainTest"}

		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		if test.AuthFlag != nil {
			arg := fmt.Sprintf("-auth=%s", *test.AuthFlag)
			os.Args = append(os.Args, arg)
		}

		if test.DebugFlag {
			os.Args = append(os.Args, "-debug")
		}

		if test.VerboseFlag {
			os.Args = append(os.Args, "-verbose")
		}

		_panicOnExit = false

		actualAuth, actualDebug, actualVerbose := GetParameters()

		assert.Equal(s.T(), test.ExpectedAuth, actualAuth, "GetParameters() Auth test '%s'", test.Description)
		assert.Equal(s.T(), test.ExpectedDebug, actualDebug, "GetParameters() Debug test '%s'", test.Description)
		assert.Equal(s.T(), test.ExpectedVerbose, actualVerbose, "GetParameters() Verbose test '%s'", test.Description)
	}
}

func (s *GithubForkUpdateSuite) Test_GetParameters_AuthFlag_Empty() {
	os.Args = []string{"mainTest"}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Args = append(os.Args, "-auth=")
	os.Args = append(os.Args, "-debug")
	os.Args = append(os.Args, "-verbose")

	_panicOnExit = true

	defer func() {
		if r := recover(); r == nil {
			s.T().Errorf("The code did not panic")
		} else {
			s.T().Logf("Recovered in %v", r)
		}
	}()

	GetParameters()

	assert.Fail(s.T(), "Test_GetParameters_AuthFlag_Empty expected Panic to fire")
}

func (s *GithubForkUpdateSuite) Test_GetParameters_FlagParse() {
	os.Args = []string{"mainTest"}

	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	os.Args = append(os.Args, "-auth=''")
	os.Args = append(os.Args, "-debug")
	os.Args = append(os.Args, "-verbose")
	os.Args = append(os.Args, "-panic")

	_panicOnExit = true

	defer func() {
		if r := recover(); r == nil {
			s.T().Errorf("The code did not panic")
		} else {
			s.T().Logf("Recovered in %v", r)
		}
	}()

	GetParameters()

	assert.Fail(s.T(), "Test_GetParameters_FlagParse expected Panic to fire")
}
