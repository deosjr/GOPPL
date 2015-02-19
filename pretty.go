
package prolog

import "fmt"

func Print_memory() {
	for k,v := range memory {
		for _,rule := range v {
			fmt.Printf("%s(", k.functor)
			for i,h := range rule.head {
				if i == k.arity-1 {
					fmt.Printf("%s)", h.Term_to_string())
					break
				}
				fmt.Printf("%s,",h.Term_to_string())
			}
			if len(rule.body) == 0 {
				fmt.Println(".")
			} else {
				fmt.Println(" :-")
				for i,b := range rule.body {
					if i == len(rule.body)-1 {
						fmt.Printf("\t%s.", b.Term_to_string())
						break
					}
					fmt.Printf("\t%s,\n",b.Term_to_string())
				}
				fmt.Println()
			}
		}
		fmt.Println()
	}
}

func Print_answer(query []Term, answer chan Alias) {
	fmt.Printf("?- %s.\n", query[0].Term_to_string())
	wait := false
	for alias := range answer {
		for k,v := range alias {
			fmt.Printf("%s = %s. ", k, v.Term_to_string())
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