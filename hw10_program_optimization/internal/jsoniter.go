//go:build jsoniter

package internal

import jsoniter "github.com/json-iterator/go"

var (
	json      = jsoniter.ConfigCompatibleWithStandardLibrary
	Unmarshal = json.Unmarshal
)
