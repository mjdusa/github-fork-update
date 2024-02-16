package environment_test

import (
	"os"
	"runtime/debug"
	"testing"

	"github.com/mjdusa/github-fork-update/internal/environment"
	"github.com/mjdusa/github-fork-update/internal/version"
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

func (s *EnvSuite) TestGetParameters2() {
	tests := []struct {
		name     string
		args     []string
		wantDbg  bool
		wantErr  bool
		wantVerb bool
	}{
		{
			name:     "Test with no arguments",
			args:     []string{},
			wantDbg:  false,
			wantErr:  true,
			wantVerb: false,
		},
		{
			name:     "Test with auth empty",
			args:     []string{"-auth", ""},
			wantDbg:  false,
			wantErr:  true,
			wantVerb: false,
		},
		{
			name:     "Test with auth argument",
			args:     []string{"-auth", "test_token"},
			wantDbg:  false,
			wantErr:  false,
			wantVerb: false,
		},
		{
			name:     "Test with debug argument",
			args:     []string{"-auth", "test_token", "-debug"},
			wantDbg:  true,
			wantErr:  false,
			wantVerb: false,
		},
		{
			name:     "Test with verbose argument",
			args:     []string{"-auth", "test_token", "-verbose"},
			wantDbg:  false,
			wantErr:  false,
			wantVerb: true,
		},
		{
			name:     "Test with both debug and verbose argument",
			args:     []string{"-auth", "test_token", "-debug", "-verbose"},
			wantDbg:  true,
			wantErr:  false,
			wantVerb: true,
		},
		{
			name:     "Test with invalid argument",
			args:     []string{"-invalid", "value"},
			wantDbg:  false,
			wantErr:  true,
			wantVerb: false,
		},
	}

	env, err := environment.NewEnvironment()
	if err != nil {
		s.T().Errorf("NewEnvironment() error = %v", err)
	}

	for _, tst := range tests {
		s.T().Run(tst.name, func(t *testing.T) {
			// Set the command line arguments
			os.Args = append([]string{"app"}, tst.args...)

			_, gotDbg, gotVerb, err := env.GetParameters()

			if err != nil {
				if !tst.wantErr {
					t.Errorf("GetParameters() test %s returned error = %v, wantErr %v", tst.name, err, tst.wantErr)
				}
			} else {
				if tst.wantErr {
					t.Errorf("GetParameters() test %s returned no error, wantErr %v", tst.name, tst.wantErr)
				}
				assert.Equal(t, tst.wantDbg, *gotDbg, "GetParameters() Debug test '%s'", tst.name)
				assert.Equal(t, tst.wantVerb, *gotVerb, "GetParameters() Verbose test '%s'", tst.name)
			}
		})
	}
}

func (s *EnvSuite) TestReport() {
	var info string

	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		info = buildInfo.String()
	}

	tests := []struct {
		name    string
		verbose bool
		dbg     bool
		want    string
	}{
		{
			name:    "Test with verbose and dbg both false",
			verbose: false,
			dbg:     false,
			want:    "",
		},
		{
			name:    "Test with verbose true and dbg false",
			verbose: true,
			dbg:     false,
			want:    version.GetVersion(),
		},
		{
			name:    "Test with verbose false and dbg true",
			verbose: false,
			dbg:     true,
			want:    info,
		},
		{
			name:    "Test with verbose and dbg both true",
			verbose: true,
			dbg:     true,
			want:    version.GetVersion() + info,
		},
	}

	env, err := environment.NewEnvironment()
	if err != nil {
		s.T().Errorf("NewEnvironment() error = %v", err)
	}

	for _, tst := range tests {
		s.T().Run(tst.name, func(t *testing.T) {
			if got := env.Report(tst.verbose, tst.dbg); got != tst.want {
				t.Errorf("Report() = %v, want %v", got, tst.want)
			}
		})
	}
}
