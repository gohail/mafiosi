package main

import (
	"encoding/json"
	"github.com/sanity-io/litter"
	"testing"
)

func Test_Main(t *testing.T) {
	s := struct {
		Text  string `json:"text"`
		Inter int    `json:"inter"`
	}{
		Text:  "Option",
		Inter: 23,
	}
	litter.Dump(s)
	body, _ := json.Marshal(s)
	litter.Dump(string(body))
}
