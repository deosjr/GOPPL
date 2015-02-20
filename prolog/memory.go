
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
	return InitPeano()
}

//TODO: move to separate testfiles!
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