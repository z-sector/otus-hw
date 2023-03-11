//go:build !go_json && !jsoniter

package internal

import "encoding/json"

var Unmarshal = json.Unmarshal
