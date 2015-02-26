package main

import (
	"io"
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	
	"GOPPL/memory"
	"GOPPL/prolog"
)

// TODO: query int(X) to Peano kicks off infinite go routines
//		 this is not very efficient
// TODO: multi-term queries: int(X), int(Y).
// 		 Right now only Y is in scope in answer

func main() {

	var file string
	flag.StringVar(&file, "f", "", ".pl or .pro file")
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
	if err == io.EOF {
		panic(reader.ThrowError(memory.ErrQueryError))
	}
	if err != nil {
		panic(err)
	}
	var_alias := make(map[prolog.VarTemplate]prolog.Term)
	query := prolog.Terms{}
	for _, t := range terms {
		qt, var_alias := t.CreateVars(var_alias)
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
		//memory.PrintMemory()	// TODO: parse listing/1
		fmt.Print("?- ")
		var input string
		fmt.Scanln(&input)
		//TODO: parse something other than query, such as exit/1
		//		or [filename] to load a file (rather consult/1)
		
		query := parseQuery(input)
		empty, answer := prolog.GetInit()
		go prolog.DFS(query, empty, answer)
		
		wait := true
		ANSWERS:
		for alias := range answer {
			if len(alias) == 0 {
				fmt.Print("True.")
			} else {
				for k,v := range alias {
					fmt.Printf("%s = %s. ", k, v.String())
				}
			}
			if wait {
				WAIT:
				for {
					var response string
					fmt.Scanln(&response)
					switch response {
					case ";": 
						break WAIT
					case "a":
						wait = false
						break WAIT
					case "q":
						break ANSWERS
					}
				}
			}
			fmt.Println()
		}
		fmt.Println("False.")
	
	}
}
