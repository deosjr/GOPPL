
package memory

import (
	"fmt"
	"os"
	
	t "GOPPL/prolog"	// TODO: dont alias, once all those Inits are gone!
)

func addData(pred t.Predicate, r t.Rule) {
	if value, ok := t.Memory[pred]; ok {
		t.Memory[pred] = append(value, r)
	} else {
		t.Memory[pred] = []t.Rule{r}
	}
}

//TODO: take a .pl file as input and parse
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

func InitMemory() t.Terms {
	InitBuiltIns()
	return InitPerms()
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

func InitPeano() t.Terms {
	s := t.Predicate{"s",1}
	addData(t.Predicate{"int",1}, t.Rule{t.Terms{t.Atom{"0"}}, t.Terms{}})
	m := &t.Var{"M"}
	sm := t.Compound_Term{s, t.Terms{m}}
	im := t.Compound_Term{t.Predicate{"int",1}, t.Terms{m}}
	addData(t.Predicate{"int",1}, t.Rule{t.Terms{sm}, t.Terms{im}})
	n := &t.Var{"N"}
	addData(t.Predicate{"sum",3}, t.Rule{t.Terms{t.Atom{"0"},m,m}, t.Terms{}})
	k := &t.Var{"K"}
	sn := t.Compound_Term{s, t.Terms{n}}
	sk := t.Compound_Term{s, t.Terms{k}}
	snmk := t.Compound_Term{t.Predicate{"sum",3}, t.Terms{n,m,k}}
	addData(t.Predicate{"sum",3}, t.Rule{t.Terms{sn, m, sk}, t.Terms{snmk}})
	
	x := &t.Var{"X"}
	//query := t.Terms{t.Compound_Term{t.Predicate{"int",1}, t.Terms{x}}}
	s2 := t.Compound_Term{s, t.Terms{t.Compound_Term{s, t.Terms{t.Atom{"0"}}}}}
	s3 := t.Compound_Term{s, t.Terms{t.Compound_Term{s, t.Terms{t.Compound_Term{s, t.Terms{t.Atom{"0"}}}}}}}
	sum := t.Compound_Term{t.Predicate{"sum",3}, t.Terms{s2,s3,x}}
	query := t.Terms{sum}
	return query
}

func InitPerms() t.Terms {
	addData(t.Predicate{"sym",1}, t.Rule{t.Terms{t.Atom{"a"}}, t.Terms{}})
	addData(t.Predicate{"sym",1}, t.Rule{t.Terms{t.Atom{"b"}}, t.Terms{}})
	h1 := &t.Var{"H1"}
	h2 := &t.Var{"H2"}
	sym1 := t.Compound_Term{t.Predicate{"sym",1}, t.Terms{h1}}
	sym2 := t.Compound_Term{t.Predicate{"sym",1}, t.Terms{h2}}
	addData(t.Predicate{"hardcoded2", 2}, t.Rule{t.Terms{h1,h2}, t.Terms{sym1, sym2}})
	x,y := &t.Var{"X"}, &t.Var{"Y"}
	h := t.Compound_Term{t.Predicate{"hardcoded2",2}, t.Terms{x,y}}
	query := t.Terms{h}
	return query
}

func InitExample() t.Terms {
	addData(t.Predicate{"p",1}, t.Rule{t.Terms{t.Atom{"a"}}, t.Terms{}})
	x := &t.Var{"X"}
	qx1 := t.Compound_Term{t.Predicate{"q",1}, t.Terms{x}}
	rx1 := t.Compound_Term{t.Predicate{"r",1}, t.Terms{x}}
	addData(t.Predicate{"p",1}, t.Rule{t.Terms{x}, t.Terms{qx1, rx1}})
	ux2 := t.Compound_Term{t.Predicate{"u",1}, t.Terms{x}}
	addData(t.Predicate{"p",1}, t.Rule{t.Terms{x}, t.Terms{ux2}})
	sx3 := t.Compound_Term{t.Predicate{"s",1}, t.Terms{x}}
	addData(t.Predicate{"q",1}, t.Rule{t.Terms{x}, t.Terms{sx3}})
	addData(t.Predicate{"r",1}, t.Rule{t.Terms{t.Atom{"a"}}, t.Terms{}})
	addData(t.Predicate{"r",1}, t.Rule{t.Terms{t.Atom{"b"}}, t.Terms{}})
	addData(t.Predicate{"s",1}, t.Rule{t.Terms{t.Atom{"a"}}, t.Terms{}})
	addData(t.Predicate{"s",1}, t.Rule{t.Terms{t.Atom{"b"}}, t.Terms{}})
	addData(t.Predicate{"s",1}, t.Rule{t.Terms{t.Atom{"c"}}, t.Terms{}})
	addData(t.Predicate{"u",1}, t.Rule{t.Terms{t.Atom{"d"}}, t.Terms{}})
	
	px := t.Compound_Term{t.Predicate{"p",1}, t.Terms{x}}
	query := t.Terms{px}
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