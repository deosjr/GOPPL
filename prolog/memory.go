
package prolog

import (
	"os"
	"GOPPL/parser"
	t "GOPPL/types"
)

var memory map[t.Predicate][]t.Rule = make(map[t.Predicate][]t.Rule)

func ruleToMem(pred t.Predicate, r t.Rule) {
	if value, ok := memory[pred]; ok {
		memory[pred] = append(value, r)
	} else {
		memory[pred] = []t.Rule{r}
	}
}

//TODO: take a .pl file as input and parse
func InitFromFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	
	reader := parser.NewReader(f)
	reader.Read()
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
	ruleToMem(t.Predicate{"UNIFY",2}, t.Rule{t.Terms{x, x}, t.Terms{}})

	// Lists as LIST/2 using t.Atom EMPTYLIST as [] and RESERVED as end of list
	// TODO: How come anon t.Vars already seem to work?!
	list := t.Predicate{"LIST",2}
	ruleToMem(list, t.Rule{t.Terms{t.Atom{"EMPTYLIST"}, t.Atom{"RESERVED"}}, t.Terms{}})
	empty_list := t.List{t.Compound_Term{list, t.Terms{t.Atom{"EMPTYLIST"}, t.Atom{"RESERVED"}}}}
	tlist := t.List{t.Compound_Term{list, t.Terms{&t.Var{"_"}, t.List{t.Compound_Term{list, t.Terms{&t.Var{"_"}, empty_list}}}}}}
	ruleToMem(list, t.Rule{t.Terms{&t.Var{"_"}, tlist}, t.Terms{}})

}

//TODO: move to separate testfiles!
func InitLists() t.Terms {

	list := t.Predicate{"LIST",2}
	empty_list := t.List{t.Compound_Term{list, t.Terms{t.Atom{"EMPTYLIST"}, t.Atom{"RESERVED"}}}}
	
	// now lets try concatenation
	l := &t.Var{"L"}
	ruleToMem(t.Predicate{"cat",3}, t.Rule{t.Terms{empty_list, l, l}, t.Terms{}})
	h, tail, r := &t.Var{"H"}, &t.Var{"T"}, &t.Var{"R"}
	ht := t.List{t.Compound_Term{list, t.Terms{h, tail}}}
	hr  := t.List{t.Compound_Term{list, t.Terms{h, r}}}
	reccat := t.Compound_Term{t.Predicate{"cat",3}, t.Terms{tail,l,r}}
	ruleToMem(t.Predicate{"cat",3}, t.Rule{t.Terms{ht, l, hr}, t.Terms{reccat}})
	
	// query
	l12345 := t.List{t.Compound_Term{list, t.Terms{t.Atom{"1"}, t.List{t.Compound_Term{list, t.Terms{t.Atom{"2"}, t.List{t.Compound_Term{list, t.Terms{t.Atom{"3"}, t.List{t.Compound_Term{list, t.Terms{t.Atom{"4"}, t.List{t.Compound_Term{list, t.Terms{t.Atom{"5"}, empty_list}}}}}}}}}}}}}}}
	x := &t.Var{"X"}
	cat := t.Compound_Term{t.Predicate{"cat",3}, t.Terms{l,x,l12345}}
	query := t.Terms{cat}
	
	return query
}

func InitPeano() t.Terms {
	s := t.Predicate{"s",1}
	ruleToMem(t.Predicate{"int",1}, t.Rule{t.Terms{t.Atom{"0"}}, t.Terms{}})
	m := &t.Var{"M"}
	sm := t.Compound_Term{s, t.Terms{m}}
	im := t.Compound_Term{t.Predicate{"int",1}, t.Terms{m}}
	ruleToMem(t.Predicate{"int",1}, t.Rule{t.Terms{sm}, t.Terms{im}})
	n := &t.Var{"N"}
	ruleToMem(t.Predicate{"sum",3}, t.Rule{t.Terms{t.Atom{"0"},m,m}, t.Terms{}})
	k := &t.Var{"K"}
	sn := t.Compound_Term{s, t.Terms{n}}
	sk := t.Compound_Term{s, t.Terms{k}}
	snmk := t.Compound_Term{t.Predicate{"sum",3}, t.Terms{n,m,k}}
	ruleToMem(t.Predicate{"sum",3}, t.Rule{t.Terms{sn, m, sk}, t.Terms{snmk}})
	
	x := &t.Var{"X"}
	//query := t.Terms{t.Compound_Term{t.Predicate{"int",1}, t.Terms{x}}}
	s2 := t.Compound_Term{s, t.Terms{t.Compound_Term{s, t.Terms{t.Atom{"0"}}}}}
	s3 := t.Compound_Term{s, t.Terms{t.Compound_Term{s, t.Terms{t.Compound_Term{s, t.Terms{t.Atom{"0"}}}}}}}
	sum := t.Compound_Term{t.Predicate{"sum",3}, t.Terms{s2,s3,x}}
	query := t.Terms{sum}
	return query
}

func InitPerms() t.Terms {
	ruleToMem(t.Predicate{"sym",1}, t.Rule{t.Terms{t.Atom{"a"}}, t.Terms{}})
	ruleToMem(t.Predicate{"sym",1}, t.Rule{t.Terms{t.Atom{"b"}}, t.Terms{}})
	h1 := &t.Var{"H1"}
	h2 := &t.Var{"H2"}
	sym1 := t.Compound_Term{t.Predicate{"sym",1}, t.Terms{h1}}
	sym2 := t.Compound_Term{t.Predicate{"sym",1}, t.Terms{h2}}
	ruleToMem(t.Predicate{"hardcoded2", 2}, t.Rule{t.Terms{h1,h2}, t.Terms{sym1, sym2}})
	x,y := &t.Var{"X"}, &t.Var{"Y"}
	h := t.Compound_Term{t.Predicate{"hardcoded2",2}, t.Terms{x,y}}
	query := t.Terms{h}
	return query
}

func InitExample() t.Terms {
	ruleToMem(t.Predicate{"p",1}, t.Rule{t.Terms{t.Atom{"a"}}, t.Terms{}})
	x := &t.Var{"X"}
	qx1 := t.Compound_Term{t.Predicate{"q",1}, t.Terms{x}}
	rx1 := t.Compound_Term{t.Predicate{"r",1}, t.Terms{x}}
	ruleToMem(t.Predicate{"p",1}, t.Rule{t.Terms{x}, t.Terms{qx1, rx1}})
	ux2 := t.Compound_Term{t.Predicate{"u",1}, t.Terms{x}}
	ruleToMem(t.Predicate{"p",1}, t.Rule{t.Terms{x}, t.Terms{ux2}})
	sx3 := t.Compound_Term{t.Predicate{"s",1}, t.Terms{x}}
	ruleToMem(t.Predicate{"q",1}, t.Rule{t.Terms{x}, t.Terms{sx3}})
	ruleToMem(t.Predicate{"r",1}, t.Rule{t.Terms{t.Atom{"a"}}, t.Terms{}})
	ruleToMem(t.Predicate{"r",1}, t.Rule{t.Terms{t.Atom{"b"}}, t.Terms{}})
	ruleToMem(t.Predicate{"s",1}, t.Rule{t.Terms{t.Atom{"a"}}, t.Terms{}})
	ruleToMem(t.Predicate{"s",1}, t.Rule{t.Terms{t.Atom{"b"}}, t.Terms{}})
	ruleToMem(t.Predicate{"s",1}, t.Rule{t.Terms{t.Atom{"c"}}, t.Terms{}})
	ruleToMem(t.Predicate{"u",1}, t.Rule{t.Terms{t.Atom{"d"}}, t.Terms{}})
	
	px := t.Compound_Term{t.Predicate{"p",1}, t.Terms{x}}
	query := t.Terms{px}
	return query
}