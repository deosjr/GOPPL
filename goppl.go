package main

import (
	"strings"
	
	"GOPPL/memory"
	"GOPPL/prolog"
	
)

func main() {

	file := "example.pl"
	memory.InitFromFile(file)
	memory.InitBuiltIns()
	
	// TODO: query int(X) at peano kicks off infinite go routines
	//		 this is not very efficient
	// TODO: multi-term queries: int(X), int(Y).
	// 		 Right now only Y is in scope in answer
	query := parseQuery("int(X).")

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
