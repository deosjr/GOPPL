
package memory

import (
	"fmt"
	"os"
	
	"GOPPL/prolog"
)

//TODO: occurs check, don't allow doubles!
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
			addData(pred, rule)
		}
	}
}

//TODO: suppress these by default when printing memory 
func InitBuiltIns() {

	x := prolog.VarTemplate{"X"}
	anon := prolog.VarTemplate{"_"}

	//TODO:
	//	not/1
	//	is/2 as IS
	//	\=/2 as not(UNIFY)
	
	//	=/2 as UNIFY(X,X)
	addData(prolog.Predicate{"UNIFY",2}, prolog.Rule{prolog.Terms{x, x}, prolog.Terms{}})

	// Lists as LIST/2 using prolog.Atom EMPTYLIST as [] and RESERVED as end of list
	// TODO: How come anon vars already seem to work?!
	//		 They don't, you don't ask for a variable list, just use to check
	list := prolog.Predicate{"LIST",2}
	
	// LIST([], RESERVED)
	addData(list, prolog.Rule{prolog.Terms{prolog.Atom{"EMPTYLIST"}, prolog.Atom{"RESERVED"}}, prolog.Terms{}})
	
	// LIST(_, LIST(_,_))
	tlist := prolog.CreateList(prolog.Terms{anon, anon}, prolog.Empty_List)
	addData(list, prolog.Rule{prolog.Terms{anon, tlist}, prolog.Terms{}})

}

func PrintMemory() {
	for k,v := range prolog.Memory {
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