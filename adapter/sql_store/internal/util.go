package internal

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
)

func CalSha1(s string) string {
	h := sha1.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Json(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
