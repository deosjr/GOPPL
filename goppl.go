package main

import (
	"bufio"
	"io"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	
	"GOPPL/memory"
	"GOPPL/prolog"
)

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
		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		input := s.Text()
		//TODO: parse something other than query, such as exit/1
		//		or [filename] to load a file (rather consult/1)
		
		query := parseQuery(input)
		node := prolog.StartDFS(query)
		
		wait := true
		s.Split(bufio.ScanRunes)
		ANSWERS:
		for result := range node.Answer {
			alias := result.A
			if result.Err == prolog.Notification {
				break
			}
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
					s.Scan()
					switch s.Text() {
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
			node.Notify()
			if !wait {
				fmt.Println()
			}
		}
		fmt.Println("False.")
	
	}
}
