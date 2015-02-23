
package main

import (
	"GOPPL/memory"
	"GOPPL/prolog"
)

func main() {	

	//file := "example.pl"
	file := "tests/permutation_test.pl"
	memory.InitFromFile(file)

	query := memory.InitPerms()
	
	memory.PrintMemory()
	
	stack := prolog.InitStack(query)
	answer := make(chan prolog.Alias, 1)
	go prolog.DFS(stack, answer)
	
	prolog.PrintAnswer(query, answer)

}