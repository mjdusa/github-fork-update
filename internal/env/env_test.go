package env_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/mjdusa/github-fork-update/internal/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Setup Suite
type EnvSuite struct {
	suite.Suite
}

func Test_EnvSuite(t *testing.T) {
	suite.Run(t, &EnvSuite{})
}

type TestGetParameters struct {
	Description     string
	AuthFlag        *string
	DebugFlag       *bool
	VerboseFlag     *bool
	ExpectedError   error
	ExpectedAuth    *string
	ExpectedDebug   *bool
	ExpectedVerbose *bool
}

func (s *EnvSuite) Test_GetParameters() {
	emptyAuth := ""
	expectedAuth := "foo-bar"
	expectedFalse := false
	expectedTrue := true

	testList := []TestGetParameters{
		{
			Description:     "Has nil Token value",
			AuthFlag:        nil,
			DebugFlag:       &expectedFalse,
			VerboseFlag:     &expectedFalse,
			ExpectedError:   fmt.Errorf("error missing auth token"),
			ExpectedAuth:    nil,
			ExpectedDebug:   &expectedFalse,
			ExpectedVerbose: &expectedFalse,
		},
		{
			Description:     "Has empty Token value",
			AuthFlag:        &emptyAuth,
			DebugFlag:       &expectedFalse,
			VerboseFlag:     &expectedFalse,
			ExpectedError:   fmt.Errorf("error missing auth token"),
			ExpectedAuth:    &emptyAuth,
			ExpectedDebug:   &expectedFalse,
			ExpectedVerbose: &expectedFalse,
		},
		{
			Description:     "Has only Token value",
			AuthFlag:        &expectedAuth,
			DebugFlag:       &expectedFalse,
			VerboseFlag:     &expectedFalse,
			ExpectedError:   nil,
			ExpectedAuth:    &expectedAuth,
			ExpectedDebug:   &expectedFalse,
			ExpectedVerbose: &expectedFalse,
		},
		{
			Description:     "Has all values, Debug value false",
			AuthFlag:        &expectedAuth,
			DebugFlag:       &expectedFalse,
			VerboseFlag:     &expectedFalse,
			ExpectedError:   nil,
			ExpectedAuth:    &expectedAuth,
			ExpectedDebug:   &expectedFalse,
			ExpectedVerbose: &expectedFalse,
		},
		{
			Description:     "Has all values, Verbose value true",
			AuthFlag:        &expectedAuth,
			DebugFlag:       &expectedFalse,
			VerboseFlag:     &expectedTrue,
			ExpectedError:   nil,
			ExpectedAuth:    &expectedAuth,
			ExpectedDebug:   &expectedFalse,
			ExpectedVerbose: &expectedTrue,
		},
		{
			Description:     "Has all values, Debug and Verbose value true",
			AuthFlag:        &expectedAuth,
			DebugFlag:       &expectedTrue,
			VerboseFlag:     &expectedTrue,
			ExpectedError:   nil,
			ExpectedAuth:    &expectedAuth,
			ExpectedDebug:   &expectedTrue,
			ExpectedVerbose: &expectedTrue,
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

		actualAuth, actualDebug, actualVerbose, actualError := env.GetParameters()

		if actualError == nil {
			if test.ExpectedAuth == nil {
				assert.Nil(s.T(), actualAuth, "GetParameters() Auth test '%s'", test.Description)
			} else {
				assert.Equal(s.T(), *test.ExpectedAuth, *actualAuth, "GetParameters() Auth test '%s'", test.Description)
			}
			assert.Equal(s.T(), *test.ExpectedDebug, *actualDebug, "GetParameters() Debug test '%s'", test.Description)
			assert.Equal(s.T(), *test.ExpectedVerbose, *actualVerbose, "GetParameters() Verbose test '%s'", test.Description)
		} else {
			assert.Equal(s.T(), fmt.Sprint(test.ExpectedError), fmt.Sprint(actualError), "GetParameters() Error test '%s'", test.Description)
		}
	}
}
