package main

import (
	"time"
)

var modules []*Module

func InitModules() {
	modules = []*Module{Serial, Time}
}

type ReturnValue struct {
	HasVal bool
	Val    string
	IsErr  bool
}

var Serial = &Module{
	Name: "serial",
	ProvidedMethods: []Method{
		{Name: "send", Fun: func(args []*Node) ReturnValue {
			println("running")
			for _, n := range args {
				println(n.Data)
				Run(n)
				print(n.Data)
			}
			return ReturnValue{HasVal: false}
		}},
	},
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

var Time = &Module{
	Name: "time",
	ProvidedMethods: []Method{
		{Name: "hold", Fun: func(args []*Node) ReturnValue {
			if len(args) != 1 {
				return Error("Invalid number of arguments provided to hold call")
			}
			Run(args[0])
			time.Sleep(GetDuration(args[0].Data))
			return ReturnValue{HasVal: false}
		}},
		{Name: "us", Fun: func(args []*Node) ReturnValue {
			if len(args) != 1 {
				return Error("Invalid number of arguments provided to time call")
			}
			Run(args[0])
			return ReturnValue{HasVal: true, Val: args[0].Data + "us"}
		}},
		{Name: "ms", Fun: func(args []*Node) ReturnValue {
			if len(args) != 1 {
				return Error("Invalid number of arguments provided to time call")
			}
			Run(args[0])
			return ReturnValue{HasVal: true, Val: args[0].Data + "ms"}
		}},
		{Name: "s", Fun: func(args []*Node) ReturnValue {
			if len(args) != 1 {
				return Error("Invalid number of arguments provided to time call")
			}
			Run(args[0])
			return ReturnValue{HasVal: true, Val: args[0].Data + "s"}
		}},
		{Name: "m", Fun: func(args []*Node) ReturnValue {
			if len(args) != 1 {
				return Error("Invalid number of arguments provided to time call")
			}
			Run(args[0])
			return ReturnValue{HasVal: true, Val: args[0].Data + "m"}
		}},
		{Name: "h", Fun: func(args []*Node) ReturnValue {
			if len(args) != 1 {
				return Error("Invalid number of arguments provided to time call")
			}
			Run(args[0])
			return ReturnValue{HasVal: true, Val: args[0].Data + "h"}
		}},
	},
}

type Module struct {
	Name            string
	ProvidedMethods []Method
}

type Method struct {
	Fun  func(args []*Node) ReturnValue
	Name string
}
