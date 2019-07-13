package unifi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUAPPoints(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	// We're just making sure an empty dataset does not crash the method.
	// https://github.com/davidnewhall/unifi-poller/issues/82
	u := &UAP{}
	pts, err := u.Points()
	a.Nil(err)
	a.NotNil(pts)
}
