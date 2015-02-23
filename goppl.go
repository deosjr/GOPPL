
package main

import (
	"GOPPL/prolog"
	"GOPPL/types"
)

func main() {	

	file := "example.pl"
	prolog.InitFromFile(file)

	query := prolog.InitMemory()
	
	prolog.Print_memory()
	
	stack := prolog.InitStack(query)
	answer := make(chan types.Alias, 1)
	go prolog.DFS(stack, answer)
	
	prolog.Print_answer(query, answer)

}