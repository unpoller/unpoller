package unidev

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlexInt(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	five := []byte(`{"channel": "5"}`)
	seven := []byte(`{"channel": 7}`)
	type reply struct {
		Channel FlexInt `json:"channel"`
	}
	var r reply
	a.Nil(json.Unmarshal(five, &r))
	a.EqualValues(FlexInt(5), r.Channel)
	a.Nil(json.Unmarshal(seven, &r))
	a.EqualValues(FlexInt(7), r.Channel)
}
