package util

import (
	"math/rand"
	"net"
	"strings"
	"time"
)

func RandomSleep(duration time.Duration, minRatio, maxRatio float64) time.Duration {
	min := int64(minRatio * 1000.0)
	max := int64(maxRatio * 1000.0)
	var n int64
	if max <= min {
		n = min
	} else {
		n = rand.Int63n(max-min) + min
	}
	d := duration * time.Duration(n) / time.Duration(1000)
	time.Sleep(d)
	return d
}

func IsIPv4(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && strings.Contains(ipAddr, ".")
}

func IsIPv6(ipAddr string) bool {
	ip := net.ParseIP(ipAddr)
	return ip != nil && strings.Contains(ipAddr, ":")
}
