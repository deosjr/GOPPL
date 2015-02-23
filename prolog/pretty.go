
package prolog

import (
	"fmt"
	
	"GOPPL/types"
)

func Print_memory() {
	for k,v := range memory {
		for _,rule := range v {
			fmt.Printf("%s(", k.Functor)
			for i,h := range rule.Head {
				if i == k.Arity-1 {
					fmt.Printf("%s)", h.String())
					break
				}
				fmt.Printf("%s,",h.String())
			}
			if len(rule.Body) == 0 {
				fmt.Println(".")
			} else {
				fmt.Println(" :-")
				for i,b := range rule.Body {
					if i == len(rule.Body)-1 {
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
func Print_answer(query types.Terms, answer chan types.Alias) {
	fmt.Printf("?- %s.\n", query[0].String())
	wait := true//false
	for alias := range answer {
		for k,v := range alias {
			fmt.Printf("%s = %s. ", k, v.String())
		}
		if wait {
			for {
				var response string
				fmt.Scanln(&response)
				if response == ";" { 
					break 
				}
				if response == "a" { 
					wait = false
					break 
				}
			}
		}
		fmt.Println()
	}
	fmt.Println("False.")
}