
package memory

import (
	"fmt"
	"os"
	
	t "GOPPL/prolog"	// TODO: dont alias, once all those Inits are gone!
)

//TODO: occurs check, don't allow doubles!
func addData(pred t.Predicate, r t.Rule) {
	if value, ok := t.Memory[pred]; ok {
		t.Memory[pred] = append(value, r)
	} else {
		t.Memory[pred] = []t.Rule{r}
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

	x := &t.Var{"X"}

	//TODO:
	//	not/1
	//	is/2 as IS
	//	\=/2 as not(UNIFY)
	
	//	=/2 as UNIFY
	addData(t.Predicate{"UNIFY",2}, t.Rule{t.Terms{x, x}, t.Terms{}})

	// Lists as LIST/2 using t.Atom EMPTYLIST as [] and RESERVED as end of list
	// TODO: How come anon t.Vars already seem to work?!
	list := t.Predicate{"LIST",2}
	addData(list, t.Rule{t.Terms{t.Atom{"EMPTYLIST"}, t.Atom{"RESERVED"}}, t.Terms{}})
	empty_list := t.List{t.Compound_Term{list, t.Terms{t.Atom{"EMPTYLIST"}, t.Atom{"RESERVED"}}}}
	tlist := t.List{t.Compound_Term{list, t.Terms{&t.Var{"_"}, t.List{t.Compound_Term{list, t.Terms{&t.Var{"_"}, empty_list}}}}}}
	addData(list, t.Rule{t.Terms{&t.Var{"_"}, tlist}, t.Terms{}})

}

//TODO: move to separate testfiles!
func InitLists() t.Terms {

	list := t.Predicate{"LIST",2}
	empty_list := t.List{t.Compound_Term{list, t.Terms{t.Atom{"EMPTYLIST"}, t.Atom{"RESERVED"}}}}
	
	// now lets try concatenation
	l := &t.Var{"L"}
	addData(t.Predicate{"cat",3}, t.Rule{t.Terms{empty_list, l, l}, t.Terms{}})
	h, tail, r := &t.Var{"H"}, &t.Var{"T"}, &t.Var{"R"}
	ht := t.List{t.Compound_Term{list, t.Terms{h, tail}}}
	hr  := t.List{t.Compound_Term{list, t.Terms{h, r}}}
	reccat := t.Compound_Term{t.Predicate{"cat",3}, t.Terms{tail,l,r}}
	addData(t.Predicate{"cat",3}, t.Rule{t.Terms{ht, l, hr}, t.Terms{reccat}})
	
	// query
	l12345 := t.List{t.Compound_Term{list, t.Terms{t.Atom{"1"}, t.List{t.Compound_Term{list, t.Terms{t.Atom{"2"}, t.List{t.Compound_Term{list, t.Terms{t.Atom{"3"}, t.List{t.Compound_Term{list, t.Terms{t.Atom{"4"}, t.List{t.Compound_Term{list, t.Terms{t.Atom{"5"}, empty_list}}}}}}}}}}}}}}}
	x := &t.Var{"X"}
	cat := t.Compound_Term{t.Predicate{"cat",3}, t.Terms{l,x,l12345}}
	query := t.Terms{cat}
	
	return query
}

func PrintMemory() {
	for k,v := range t.Memory {
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