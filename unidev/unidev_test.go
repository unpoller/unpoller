package unidev

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlexInt(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	type testReply struct {
		Five  FlexInt `json:"five"`
		Seven FlexInt `json:"seven"`
		Auto  FlexInt `json:"auto"`
	}
	var r testReply
	// test unmarshalling the custom type three times with different values.
	a.Nil(json.Unmarshal([]byte(`{"five": "5", "seven": 7, "auto": "auto"}`), &r))
	// test number in string.
	a.EqualValues(5, r.Five.Number)
	a.EqualValues("5", r.Five.String)
	// test number.
	a.EqualValues(7, r.Seven.Number)
	a.EqualValues("7", r.Seven.String)
	// test string.
	a.Nil(json.Unmarshal([]byte(`{"channel": "auto"}`), &r),
		"a regular string must not produce an unmarshal error")
	a.EqualValues(0, r.Auto.Number)
	a.EqualValues("auto", r.Auto.String)
}
