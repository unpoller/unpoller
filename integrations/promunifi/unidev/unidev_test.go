package unidev

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlexInt(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	type reply struct {
		Channel FlexInt `json:"channel"`
	}
	var r reply
	a.Nil(json.Unmarshal([]byte(`{"channel": "5"}`), &r))
	a.EqualValues(FlexInt(5), r.Channel)
	a.Nil(json.Unmarshal([]byte(`{"channel": 7}`), &r))
	a.EqualValues(FlexInt(7), r.Channel)
	a.Nil(json.Unmarshal([]byte(`{"channel": "auto"}`), &r),
		"a regular string must not produce an unmarshal error")
	a.EqualValues(FlexInt(0), r.Channel)
}
