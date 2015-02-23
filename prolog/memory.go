
package prolog

var memory map[Predicate][]Rule = make(map[Predicate][]Rule)

func ruleToMem(pred Predicate, r Rule) {
	if value, ok := memory[pred]; ok {
		memory[pred] = append(value, r)
	} else {
		memory[pred] = []Rule{r}
	}
}

//TODO: take a .pl file as input and parse
func InitMemory() Terms {
	InitBuiltIns()
	return InitLists()
}

//TODO: suppress these by default when printing memory 
func InitBuiltIns() {

	// atm this is the builtin definition for prolog lists
	// TODO: How come anon vars already seem to work?!
	list := Predicate{"LIST",2}
	ruleToMem(list, Rule{Terms{Atom{"EMPTYLIST"}, Atom{"RESERVED"}}, Terms{}})
	empty_list := List{Compound_Term{list, Terms{Atom{"EMPTYLIST"}, Atom{"RESERVED"}}}}
	tlist := List{Compound_Term{list, Terms{&Var{"_"}, List{Compound_Term{list, Terms{&Var{"_"}, empty_list}}}}}}
	ruleToMem(list, Rule{Terms{&Var{"_"}, tlist}, Terms{}})

}

//TODO: move to separate testfiles!
func InitLists() Terms {

	list := Predicate{"LIST",2}
	empty_list := List{Compound_Term{list, Terms{Atom{"EMPTYLIST"}, Atom{"RESERVED"}}}}
	
	// now lets try concatenation
	l := &Var{"L"}
	ruleToMem(Predicate{"cat",3}, Rule{Terms{empty_list, l, l}, Terms{}})
	h, t, r := &Var{"H"}, &Var{"T"}, &Var{"R"}
	ht := List{Compound_Term{list, Terms{h, t}}}
	hr  := List{Compound_Term{list, Terms{h, r}}}
	reccat := Compound_Term{Predicate{"cat",3}, Terms{t,l,r}}
	ruleToMem(Predicate{"cat",3}, Rule{Terms{ht, l, hr}, Terms{reccat}})
	
	// query
	l12345 := List{Compound_Term{list, Terms{Atom{"1"}, List{Compound_Term{list, Terms{Atom{"2"}, List{Compound_Term{list, Terms{Atom{"3"}, List{Compound_Term{list, Terms{Atom{"4"}, List{Compound_Term{list, Terms{Atom{"5"}, empty_list}}}}}}}}}}}}}}}
	x := &Var{"X"}
	cat := Compound_Term{Predicate{"cat",3}, Terms{l,x,l12345}}
	query := Terms{cat}
	
	return query
}

func InitPeano() Terms {
	s := Predicate{"s",1}
	ruleToMem(Predicate{"int",1}, Rule{Terms{Atom{"0"}}, Terms{}})
	m := &Var{"M"}
	sm := Compound_Term{s, Terms{m}}
	im := Compound_Term{Predicate{"int",1}, Terms{m}}
	ruleToMem(Predicate{"int",1}, Rule{Terms{sm}, Terms{im}})
	n := &Var{"N"}
	ruleToMem(Predicate{"sum",3}, Rule{Terms{Atom{"0"},m,m}, Terms{}})
	k := &Var{"K"}
	sn := Compound_Term{s, Terms{n}}
	sk := Compound_Term{s, Terms{k}}
	snmk := Compound_Term{Predicate{"sum",3}, Terms{n,m,k}}
	ruleToMem(Predicate{"sum",3}, Rule{Terms{sn, m, sk}, Terms{snmk}})
	
	x := &Var{"X"}
	//query := Terms{Compound_Term{Predicate{"int",1}, Terms{x}}}
	s2 := Compound_Term{s, Terms{Compound_Term{s, Terms{Atom{"0"}}}}}
	s3 := Compound_Term{s, Terms{Compound_Term{s, Terms{Compound_Term{s, Terms{Atom{"0"}}}}}}}
	sum := Compound_Term{Predicate{"sum",3}, Terms{s2,s3,x}}
	query := Terms{sum}
	return query
}

func InitPerms() Terms {
	ruleToMem(Predicate{"sym",1}, Rule{Terms{Atom{"a"}}, Terms{}})
	ruleToMem(Predicate{"sym",1}, Rule{Terms{Atom{"b"}}, Terms{}})
	h1 := &Var{"H1"}
	h2 := &Var{"H2"}
	sym1 := Compound_Term{Predicate{"sym",1}, Terms{h1}}
	sym2 := Compound_Term{Predicate{"sym",1}, Terms{h2}}
	ruleToMem(Predicate{"hardcoded2", 2}, Rule{Terms{h1,h2}, Terms{sym1, sym2}})
	x,y := &Var{"X"}, &Var{"Y"}
	h := Compound_Term{Predicate{"hardcoded2",2}, Terms{x,y}}
	query := Terms{h}
	return query
}

func InitExample() Terms {
	ruleToMem(Predicate{"p",1}, Rule{Terms{Atom{"a"}}, Terms{}})
	x := &Var{"X"}
	qx1 := Compound_Term{Predicate{"q",1}, Terms{x}}
	rx1 := Compound_Term{Predicate{"r",1}, Terms{x}}
	ruleToMem(Predicate{"p",1}, Rule{Terms{x}, Terms{qx1, rx1}})
	ux2 := Compound_Term{Predicate{"u",1}, Terms{x}}
	ruleToMem(Predicate{"p",1}, Rule{Terms{x}, Terms{ux2}})
	sx3 := Compound_Term{Predicate{"s",1}, Terms{x}}
	ruleToMem(Predicate{"q",1}, Rule{Terms{x}, Terms{sx3}})
	ruleToMem(Predicate{"r",1}, Rule{Terms{Atom{"a"}}, Terms{}})
	ruleToMem(Predicate{"r",1}, Rule{Terms{Atom{"b"}}, Terms{}})
	ruleToMem(Predicate{"s",1}, Rule{Terms{Atom{"a"}}, Terms{}})
	ruleToMem(Predicate{"s",1}, Rule{Terms{Atom{"b"}}, Terms{}})
	ruleToMem(Predicate{"s",1}, Rule{Terms{Atom{"c"}}, Terms{}})
	ruleToMem(Predicate{"u",1}, Rule{Terms{Atom{"d"}}, Terms{}})
	
	px := Compound_Term{Predicate{"p",1}, Terms{x}}
	query := Terms{px}
	return query
}