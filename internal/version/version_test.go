package version_test

import (
	"fmt"
	"os"
	"testing"

	sut "github.com/mjdusa/github-fork-update/internal/version" // sut - system under test.
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// Setup Suite.
type VersionSuite struct {
	suite.Suite
}

func TestVersionSuite(t *testing.T) {
	suite.Run(t, &VersionSuite{})
}

type testGetVersion struct {
	Description string
	AppVersion  string
	Branch      string
	BuildTime   string
	Commit      string
	GoVersion   string
	Expected    string
}

func testGetVersionData() []testGetVersion {
	tests := []testGetVersion{
		{
			Description: "All are empty strings",
			AppVersion:  "",
			Branch:      "",
			BuildTime:   "",
			Commit:      "",
			GoVersion:   "",
		},
		{
			Description: "All have values",
			AppVersion:  "AppVersion",
			Branch:      "Branch",
			BuildTime:   "BuildTime",
			Commit:      "Commit",
			GoVersion:   "GoVersion",
		},
	}

	return tests
}

func (s *VersionSuite) TestGetVersion() {
	for _, tst := range testGetVersionData() {
		expected := fmt.Sprintf(
			"%s version: [%s]\n- Branch:     [%s]\n- Build Time: [%s]\n- Commit:     [%s]\n- Go Version: [%s]\n",
			os.Args[0], tst.AppVersion, tst.Branch, tst.BuildTime, tst.Commit, tst.GoVersion)

		sut.AppVersion = tst.AppVersion
		sut.Branch = tst.Branch
		sut.BuildTime = tst.BuildTime
		sut.Commit = tst.Commit
		sut.GoVersion = tst.GoVersion

		actual := sut.GetVersion()

		assert.Equal(s.T(), expected, actual,
			tst.Description+fmt.Sprintf(" expected '%s', actual '%s'", expected, actual))
	}
}
