package run_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// Setup Suite
type GithubForkUpdateSuite struct {
	suite.Suite
}

func Test_GithubForkUpdate_Suite(t *testing.T) {
	suite.Run(t, &GithubForkUpdateSuite{})
}
