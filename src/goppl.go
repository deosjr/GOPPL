package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"memory"
	"prolog"
)

// TODO: multi-term queries: int(X), int(Y).
// 		 Right now only X is in scope in answer

func main() {

	var file string
	flag.StringVar(&file, "f", "", ".pl or .pro file")
	flag.Parse()

	if file == "" {
		memory.InitBuiltIns()
		REPL()
		return
	}

	if filepath.Ext(file) != ".pl" && filepath.Ext(file) != ".pro" {
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
	for i := len(terms) - 1; i >= 0; i-- {
		compound, err := reader.AtomToPredicate(terms[i])
		if err != nil {
			panic(err)
		}
		qt, var_alias := compound.CreateVars(var_alias)
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
		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		input := s.Text()

		query := parseQuery(input)
		node := prolog.EmptyDFS(query)

		wait := true
		alias := node.GetAnswer()
	ANSWERS:
		for alias != nil {
			if len(alias) == 0 {
				fmt.Print("True.")
			} else {
				for k, v := range alias {
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
			if !wait {
				fmt.Println()
			}
			alias = node.GetAnswer()
		}
		fmt.Println("False.")

	}
}
