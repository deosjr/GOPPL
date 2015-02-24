
package main

import (
	"strings"
	
	"fmt"

	"GOPPL/memory"
	"GOPPL/prolog"
)

func main() {	

	file := "example.pl"
	memory.InitFromFile(file)
	memory.InitBuiltIns()

	// TODO: query int(X) to peano_test has broken!
	query := parseQuery("int(X).")
	fmt.Println(query)
	
	memory.PrintMemory()
	
	stack := prolog.InitStack(query)
	answer := make(chan prolog.Alias, 1)
	go prolog.DFS(stack, answer)
	
	prolog.PrintAnswer(query, answer)

}

func parseQuery(q string) prolog.Terms {
	
	s := strings.NewReader(q)
	reader := memory.NewReader(s)
	
	terms, err := reader.ReadTerms()
	if err != nil {
		panic(err)
	}
	var_alias := make(map[prolog.VarTemplate]prolog.Term)
	query := prolog.Terms{}
	for _, t := range terms {
		qt, var_alias := prolog.CreateVars(t, var_alias)
		var_alias = var_alias
		query = append(query, qt)
	}
	if reader.Last_Read != reader.Stop {
		panic(reader.ThrowError(memory.ErrQueryError))
	}
	return query
	
}