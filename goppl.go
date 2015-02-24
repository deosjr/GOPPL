package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	
	"GOPPL/memory"
	"GOPPL/prolog"
)

// TODO: query int(X) at peano kicks off infinite go routines
//		 this is not very efficient
// TODO: multi-term queries: int(X), int(Y).
// 		 Right now only Y is in scope in answer

func main() {

	var file string
	flag.StringVar(&file, "f", "", "-f=.pl or .pro file")
	flag.Parse()
	
	if file == "" {
		memory.InitBuiltIns()
		REPL()
		return
	}
	
	if filepath.Ext(file) != ".pl" && filepath.Ext(file) != ".pro"{
		fmt.Println("Input a valid Prolog filename.")
		return
	}
	
	memory.InitBuiltIns()
	memory.InitFromFile(file)
	REPL()

}

func parseQuery(q string) prolog.Terms {

	s := strings.NewReader(q)
	reader := memory.NewReader(s)

	terms, err := reader.ReadTerms()
	// TODO: recover from syntax error in query (err == io.EOF)
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

func REPL() {

	for {
		fmt.Print("?- ")
		var input string
		fmt.Scanln(&input)
		//TODO: parse something other than query, such as exit/1
		//		or [filename] to load a file
		//memory.PrintMemory()	// TODO: parse listing/1
		
		query := parseQuery(input)
		stack := prolog.InitStack(query)
		answer := make(chan prolog.Alias, 1)
		go prolog.DFS(stack, answer)
		
		wait := true
		WAIT:
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
					if response == "q" {
						break WAIT
					}
				}
			}
			fmt.Println()
		}
		fmt.Println("False.")
	
	}
}
