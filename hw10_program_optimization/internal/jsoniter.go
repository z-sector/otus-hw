//go:build !go_json && !json

package internal

import jsoniter "github.com/json-iterator/go"

var (
	json      = jsoniter.ConfigCompatibleWithStandardLibrary
	Unmarshal = json.Unmarshal
)
