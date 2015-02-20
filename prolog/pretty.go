
package prolog

import "fmt"

func (a Atom) String() string{ return a.value}

func (v *Var) String() string { return v.name }

func (c Compound_Term) String() string{ 
	s := c.pred.functor + "("
	for i,t := range c.args {
		if i == c.pred.arity-1 {
			s += t.String()
			break
		}
		s += t.String() + ","
	}
	return s + ")"
}

func (tlist Terms) String() string {
	s := "["
	for _,t := range tlist {
		s = s + t.String() + " "
	}
	return s + "]"
} 

func (a Alias) String() string {
	s := "{"
	for k,v := range a {
		s = s + k.String() + ":" + v.String() + " "
	}
	return s + "}"
}

func Print_memory() {
	for k,v := range memory {
		for _,rule := range v {
			fmt.Printf("%s(", k.functor)
			for i,h := range rule.head {
				if i == k.arity-1 {
					fmt.Printf("%s)", h.String())
					break
				}
				fmt.Printf("%s,",h.String())
			}
			if len(rule.body) == 0 {
				fmt.Println(".")
			} else {
				fmt.Println(" :-")
				for i,b := range rule.body {
					if i == len(rule.body)-1 {
						fmt.Printf("\t%s.", b.String())
						break
					}
					fmt.Printf("\t%s,\n",b.String())
				}
				fmt.Println()
			}
		}
		fmt.Println()
	}
}

// Contains the ; wait loop. Set wait=false for auto all evaluations
func Print_answer(query []Term, answer chan Alias) {
	fmt.Printf("?- %s.\n", query[0].String())
	wait := false
	for alias := range answer {
		for k,v := range alias {
			fmt.Printf("%s = %s. ", k, v.String())
		}
		if wait {
			for {
				var response string
				fmt.Scanln(&response)
				if response == ";" { break }
				if response == "a" { wait=false; break }
			}
		}
		fmt.Println()
	}
	fmt.Println("False.")
}