package unifi

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlexInt(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	five, seven := 5, 7

	var r struct {
		Five    FlexInt `json:"five"`
		Seven   FlexInt `json:"seven"`
		Auto    FlexInt `json:"auto"`
		Channel FlexInt `json:"channel"`
		Nil     FlexInt `json:"nil"`
	}

	// test unmarshalling the custom type three times with different values.
	a.Nil(json.Unmarshal([]byte(`{"five": "5", "seven": 7, "auto": "auto", "nil": null}`), &r))
	// test number in string.
	a.EqualValues(five, r.Five.Val)
	a.EqualValues("5", r.Five.Txt)
	// test number.
	a.EqualValues(seven, r.Seven.Val)
	a.EqualValues("7", r.Seven.Txt)
	// test string.
	a.EqualValues(0, r.Auto.Val)
	a.EqualValues("auto", r.Auto.Txt)
	// test (error) struct.
	a.NotNil(json.Unmarshal([]byte(`{"channel": {}}`), &r),
		"a non-string and non-number must produce an error.")
	a.EqualValues(0, r.Channel.Val)
	// test null.
	a.EqualValues(0, r.Nil.Val)
	a.EqualValues("0", r.Nil.Txt)
}
