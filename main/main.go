
package main

import (
	p "GOPPL"
)

func main() {	

	query := p.InitMemory()
	
	p.Print_memory()
	
	stack := p.InitStack(query)
	answer := make(chan p.Alias, 1)
	go p.DFS(stack, answer)
	
	p.Print_answer(query, answer)

}