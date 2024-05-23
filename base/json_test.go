package base

import (
	"encoding/json"
	"testing"
)

/**
@file:
@author: levi.Tang
@time: 2024/4/17 16:41
@description:
**/

type extraParam struct {
	BootType int `json:"bootType"`
}

func TestJsonMarshal(t *testing.T) {
	str := ""
	var exts = extraParam{}
	if err := json.Unmarshal([]byte(str), &exts); err != nil {
		t.Logf("CreateGameRecord, extraParam unMarshal err: %s ", err)
		return
	}

	t.Logf("%d", exts.BootType)
}
