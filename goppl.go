
package main

import (
	"GOPPL/memory"
	"GOPPL/prolog"
)

func main() {	

	//file := "example.pl"
	file := "tests/lists_test.pl"
	memory.InitFromFile(file)
	memory.InitBuiltIns()

	query := memory.InitLists()
	
	memory.PrintMemory()
	
	stack := prolog.InitStack(query)
	answer := make(chan prolog.Alias, 1)
	go prolog.DFS(stack, answer)
	
	prolog.PrintAnswer(query, answer)

}