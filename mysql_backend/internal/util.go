package internal

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strconv"
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

func StringListToInt(ids []string) []int64 {
	return MapNoneEmpty(ids, func(item string) int64 {
		id, _ := strconv.ParseInt(item, 10, 64)
		return id
	})
}

func Map[T, F any](list []T, f func(item T) F) []F {
	res := make([]F, 0, len(list))
	for _, item := range list {
		res = append(res, f(item))
	}
	return res
}

func MapNoneEmpty[T any, F comparable](list []T, f func(item T) F) []F {
	var empty F
	res := make([]F, 0, len(list))
	for _, item := range list {
		tmp := f(item)
		if tmp == empty {
			continue
		}
		res = append(res, tmp)
	}
	return res
}

func Unique[T comparable](list []T) []T {
	if len(list) == 0 {
		return list
	}
	res := make([]T, 0, len(list))
	res = append(res, list[0])
	for i := 1; i < len(list); i++ {
		if list[i] != list[i-1] {
			res = append(res, list[i])
		}
	}
	return res
}
