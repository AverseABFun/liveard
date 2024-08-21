package main

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

func Run(n *Node) {
	if n.Data != "" && n.Parent != nil {
		Run(n.Parent)
	}
	for _, child := range n.Children {
		println(child.Data)
		switch child.Data {
		case "loop":
			if len(child.Children) > 3 {
				Error("Too many arguments to loop call(got " + Itoa(len(child.Children)) + ", expected up to 3)")
				println("debug: " + child.String())
				continue
			}
			if len(child.Children) < 2 {
				Error("Too few arguments to loop call")
				continue
			}
			var iterations, err = Atoi(child.Children[0].Data)
			if err != nil {
				Error("Expected first argument to loop call to be an int, got error: " + err.Error())
				continue
			}
			for i := 0; i < iterations; i++ {
				if child.Children[1].Data != "" {
					var broken = false
					for _, varia := range variables {
						if varia.Name == child.Children[1].Data {
							varia.Value = Itoa(i)
							broken = true
							break
						}
					}
					if !broken {
						var vari = Variable{Name: child.Children[1].Data, Value: Itoa(i)}
						variables = append(variables, &vari)
					}
					Run(child.Children[2])
				} else {
					Run(child.Children[1])
				}
			}
		default:
			var broken = false
			if child.IsVar {
				for _, vari := range variables {
					println(vari.Name)
					if vari.Name == child.Data {
						child.Data = vari.Value
						broken = true
						break
					}
				}
			}
			broken = child.IsLiteral // a little janky but it works
			if !broken {
				println("unbroken")
				for _, mod := range modules {
					println(mod.Name)
					for _, method := range mod.ProvidedMethods {
						println(method.Name)
						if method.Name == child.Data {
							println("calling")
							var out = method.Fun(child.Children)
							println("done calling")
							child.Data = out.Val
							broken = true
							println("breaking")
							break
						}
					}
				}

			}
		}
	}
}
