package main

import "time"

const DEFAULT_VARIABLE_SIZE = 100

func Error(str string) ReturnValue {
	println("Error: " + str)
	return ReturnValue{HasVal: true, IsErr: true, Val: str}
}

var variables = make([]*Variable, 0, DEFAULT_VARIABLE_SIZE)

type Variable struct {
	Name  string
	Value string
}

func deleteVar(s []*Variable, name string) []*Variable { // Copied from slices package to reduce memory footprint
	for i, val := range s {
		if val.Name == name {
			return delete(s, i, i+1)
		}
	}
	return s
}

func prepend[S ~[]E, E any](s S, e E) S {
	return append(S{e}, s...)
}

func RunSingle(child *Node) {
	args := child.Children
	switch child.Data {
	case "loop":
		if len(child.Children) > 3 {
			Error("Too many arguments to loop call(got " + Itoa(len(child.Children)) + ", expected up to 3)")
			return
		}
		if len(child.Children) < 2 {
			Error("Too few arguments to loop call")
			return
		}
		var iterations, err = Atoi(child.Children[0].Data)
		if err != nil && child.Children[0].Data != "forever" {
			Error("Expected first argument to loop call to be an int, got error: " + err.Error())
			return
		}
		if child.Children[0].Data != "forever" {
			for i := 1; i < iterations; i++ {
				if child.Children[1].Data != "" {
					variables = deleteVar(variables, child.Children[1].Data)
					var vari = Variable{Name: child.Children[1].Data, Value: Itoa(i)}
					variables = prepend(variables, &vari)
					Run(child.Children[2])
				} else {
					Run(child.Children[1])
				}
			}
		} else {
			i := 1
			for {
				if child.Children[1].Data != "" {
					variables = deleteVar(variables, child.Children[1].Data)
					var vari = Variable{Name: child.Children[1].Data, Value: Itoa(i)}
					variables = prepend(variables, &vari)
					Run(child.Children[2])
				} else {
					Run(child.Children[1])
				}
				i++
			}
		}
	case "print":
	case "out":
	case "send":
		for _, n := range args {
			Run(n)
			print(n.Data)
		}
	case "sleep":
	case "wait":
	case "hold":
		if len(args) != 1 {
			Error("Invalid number of arguments provided to hold call")
		}
		Run(args[0])
		time.Sleep(GetDuration(args[0].Data))
	case "us":
		if len(args) != 1 {
			Error("Invalid number of arguments provided to time call")
		}
		Run(args[0])
		child.Data = args[0].Data + "us"
	case "ms":
		if len(args) != 1 {
			Error("Invalid number of arguments provided to time call")
		}
		Run(args[0])
		child.Data = args[0].Data + "ms"
	case "s":
		if len(args) != 1 {
			Error("Invalid number of arguments provided to time call")
		}
		Run(args[0])
		child.Data = args[0].Data + "s"
	case "m":
		if len(args) != 1 {
			Error("Invalid number of arguments provided to time call")
		}
		Run(args[0])
		child.Data = args[0].Data + "m"
	case "h":
		if len(args) != 1 {
			Error("Invalid number of arguments provided to time call")
		}
		Run(args[0])
		child.Data = args[0].Data + "h"
	default:
		var broken = false
		if child.IsVar {
			for _, vari := range variables {
				if vari.Name == child.Data {
					child.Data = vari.Value
					broken = true
					break
				}
			}
		}
		broken = child.IsLiteral // a little janky but it works
		if !broken {
			/*for i, methodName := range methodNames {
				if methodName == child.Data {
					println("calling")
					var out = methods[i](child.Children)
					println("done calling")
					child.Data = out.Val
					broken = true
					println("breaking")
					break
				}
			}*/

		}
	}
}

func Run(n *Node) {
	if n.Data != "" && n.Parent != nil {
		RunSingle(n)
	}
	for _, child := range n.Children {
		RunSingle(child)
	}
}
