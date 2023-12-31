package version_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/mjdusa/github-fork-update/internal/version"
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

func (s *GithubForkUpdateSuite) Test_GetVersion_unpopulated() {
	app := ""
	if len(os.Args) > 0 {
		app = os.Args[0]
	}
	expected := fmt.Sprintf("%s version: []\n- Branch:     []\n- Build Time: []\n- Commit:     []\n- Go Version: []\n", app)

	version.AppVersion = ""
	version.Branch = ""
	version.BuildTime = ""
	version.Commit = ""
	version.GoVersion = ""

	actual := version.GetVersion()

	assert.Equal(s.T(), expected, actual, "GetVersion() unpopulated message expected '%s', but got '%s'", expected, actual)
}

func (s *GithubForkUpdateSuite) Test_GetVersion_populated() {
	app := ""
	if len(os.Args) > 0 {
		app = os.Args[0]
	}
	expected := fmt.Sprintf("%s version: [v1.2.3]\n- Branch:     [main]\n- Build Time: [01/01/1970T00:00:00.0000 GMT]\n- Commit:     [1234567890abcdef]\n- Go Version: [1.20.5]\n", app)

	version.AppVersion = "v1.2.3"
	version.Branch = "main"
	version.BuildTime = "01/01/1970T00:00:00.0000 GMT"
	version.Commit = "1234567890abcdef"
	version.GoVersion = "1.20.5"

	actual := version.GetVersion()

	assert.Equal(s.T(), expected, actual, "GetVersion() populated message expected '%s', but got '%s'", expected, actual)
}
