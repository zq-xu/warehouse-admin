package utils

import (
	"strconv"
	"time"
)

func OptStringPtr(dst, src *string) {
	if src == nil || dst == nil {
		return
	}

	*dst = *src
}

func OptFloat32Ptr(dst, src *float32) {
	if src == nil || dst == nil {
		return
	}

	*dst = *src
}

func OptIntPtr(dst, src *int) {
	if src == nil || dst == nil {
		return
	}

	*dst = *src
}

func OptInt64Ptr(dst, src *int64) {
	if src == nil || dst == nil {
		return
	}

	*dst = *src
}

func OptUnixTimePtr(dst, src *UnixTime) {
	if src == nil || dst == nil {
		return
	}

	*dst = *src
}

func OptTimePtrByUnixTimePtr(dst **time.Time, src *UnixTime) {
	if src == nil || dst == nil {
		return
	}

	*dst = (*time.Time)(src)
}

func OptInt64ByStringPtr(dst *int64, src *string) {
	if src == nil || dst == nil {
		return
	}

	*dst, _ = strconv.ParseInt(*src, 10, 64)
}

func GetStringFromPtr(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

func GetIntFromPtr(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

func GetFloat32FromPtr(ptr *float32) float32 {
	if ptr == nil {
		return 0
	}
	return *ptr
}
