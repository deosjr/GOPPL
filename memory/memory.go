
package memory

import (
	"fmt"
	"os"
	
	"GOPPL/prolog"
)

// TODO: occurs check, don't allow doubles!
// This should catch redefining builtin and extralogical as well
func addData(pred prolog.Predicate, r prolog.Rule) {
	if value, ok := prolog.Memory[pred]; ok {
		prolog.Memory[pred] = append(value, r)
	} else {
		prolog.Memory[pred] = []prolog.Rule{r}
	}
}

func InitFromFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	
	reader := NewReader(f)
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	
	for pred, rules := range data {
		for _, rule := range rules {
			//err = checkSingletons(rule)
			if err != nil {
				panic(err)
			}
			addData(pred, rule)
		}
	}
	//printMemory()
}

/**
func checkSingletons(rule prolog.Rule) error {
	vars := make(map[prolog.VarTemplate]int)
	for _, v := range varTemplates(append(rule.Head, rule.Body...)) {
		i, ok := vars[v]
		if ok {
			vars[v] = i+1
		} else {
			vars[v] = 1
		}
	}
	for _, v := range vars {
		if v < 2 {
			return errors.New("singleton error")
		}
	}
	return nil
}

// TODO: duplicate of search.varsInTermArgs in functionality!
func varTemplates(terms prolog.Terms) []prolog.VarTemplate {
	vars := []prolog.VarTemplate{}
	for _,term := range terms {
		switch term.(type) {
		case prolog.VarTemplate:
			vars = append(vars, term.(prolog.VarTemplate))
		case prolog.Compound_Term:
			vars = append(vars, varTemplates(term.(prolog.Compound_Term).GetArgs())...)
		}
	}
	return vars
}
*/

func printMemory() {
	RULES:
	for k,v := range prolog.Memory {
		for _,v := range builtins {
			if k == v {
				continue RULES
			}
		}
		for _,rule := range v {
			fmt.Printf("%s(", k.Functor)
			for i,h := range rule.Head {
				if i == k.Arity-1 {
					fmt.Printf("%s)", h.String())
					break
				}
				fmt.Printf("%s,",h.String())
			}
			if len(rule.Body) == 0 {
				fmt.Println(".")
			} else {
				fmt.Println(" :-")
				for i,b := range rule.Body {
					if i == len(rule.Body)-1 {
						fmt.Printf("\t%s.", b.String())
						break
					}
					fmt.Printf("\t%s,\n",b.String())
				}
				fmt.Println()
			}
		}
		fmt.Println()
	}
}