package cpanel

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/sebdah/goldie"
	"github.com/stretchr/testify/assert"
)

func TestListAllResellerNames(t *testing.T) {
	g := goldie.New(t,
		goldie.WithFixtureDir("testdata/golden"),
		goldie.WithDiffEngine(goldie.ColoredDiff),
		goldie.WithNameSuffix(".golden"),
	)

	resellers, err := testWhmApi.ListAllResellerNames()
	assert.NoError(t, err)

	g.Assert(t, t.Name(), []byte(spew.Sdump(resellers)))
}

func TestResellerUsers(t *testing.T) {
	g := goldie.New(t,
		goldie.WithFixtureDir("testdata/golden"),
		goldie.WithDiffEngine(goldie.ColoredDiff),
		goldie.WithNameSuffix(".golden"),
	)

	resellerUsers, err := testWhmApi.ResellerUsers("whousescpanel")
	assert.NoError(t, err)

	g.Assert(t, t.Name(), []byte(spew.Sdump(resellerUsers)))
}
