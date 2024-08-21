package main

import (
	"time"
)

type ReturnValue struct {
	HasVal bool
	Val    string
	IsErr  bool
}

func GetDuration(str string) time.Duration {
	if HasSuffix(str, "us") {
		var t, _ = Atoi(TrimSuffix(str, "us"))
		return time.Microsecond * time.Duration(t)
	}
	if HasSuffix(str, "ms") {
		var t, _ = Atoi(TrimSuffix(str, "ms"))
		return time.Millisecond * time.Duration(t)
	}
	if HasSuffix(str, "s") {
		var t, _ = Atoi(TrimSuffix(str, "s"))
		return time.Second * time.Duration(t)
	}
	if HasSuffix(str, "m") {
		var t, _ = Atoi(TrimSuffix(str, "m"))
		return time.Minute * time.Duration(t)
	}
	if HasSuffix(str, "h") {
		var t, _ = Atoi(TrimSuffix(str, "h"))
		return time.Hour * time.Duration(t)
	}
	return 0
}

type Method struct {
	Fun  func(args []*Node) ReturnValue
	Name string
}
