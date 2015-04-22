// Copyright 2015 unirita Inc.
// Created 2015/04/10 shanxia

package gen

import (
	"encoding/json"
)

// JSON形式のジェネレーター
type JsonGenerator struct {
}

func (s JsonGenerator) Generate(out *OutputRoot) (string, error) {
	byteMessage, err := json.Marshal(out)
	if err != nil {
		return "", err
	}
	return string(byteMessage), nil
}
