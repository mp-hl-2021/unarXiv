package utils

import "time"

func Uint64Time(timestamp time.Time) uint64 {
	return uint64(timestamp.UTC().UnixNano())
}

func TimeFromUint64(nsecs uint64) time.Time {
	return time.Unix(0, int64(nsecs)).Local()
}
