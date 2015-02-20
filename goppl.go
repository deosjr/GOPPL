
package main

import (
	"GOPPL/prolog"
)

func main() {	

	query := prolog.InitMemory()
	
	prolog.Print_memory()
	
	stack := prolog.InitStack(query)
	answer := make(chan prolog.Alias, 1)
	go prolog.DFS(stack, answer)
	
	prolog.Print_answer(query, answer)

}