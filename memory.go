
package prolog

var memory map[Predicate][]Rule = make(map[Predicate][]Rule)

func ruleToMem(pred Predicate, r Rule) {
	if value, ok := memory[pred]; ok {
		memory[pred] = append(value, r)
	} else {
		memory[pred] = []Rule{r}
	}
}

func InitMemory() []Term {
	/**
	ruleToMem(Predicate{"sym",1}, Rule{[]Term{Atom{"a"}}, []Term{}})
	ruleToMem(Predicate{"sym",1}, Rule{[]Term{Atom{"b"}}, []Term{}})
	h1 := &Var{"H1"}
	h2 := &Var{"H2"}
	sym1 := Compound_Term{Predicate{"sym",1}, []Term{h1}}
	sym2 := Compound_Term{Predicate{"sym",1}, []Term{h2}}
	ruleToMem(Predicate{"hardcoded2", 2}, Rule{[]Term{h1,h2}, []Term{sym1, sym2}})
	x,y := &Var{"X"}, &Var{"Y"}
	h := Compound_Term{Predicate{"hardcoded2",2}, []Term{x,y}}
	query := []Term{h}
	*/
	/**
	ruleToMem(Predicate{"p",1}, Rule{[]Term{Atom{"a"}}, []Term{}})
	x1 := &Var{"X"}
	qx1 := Compound_Term{Predicate{"q",1}, []Term{x1}}
	rx1 := Compound_Term{Predicate{"r",1}, []Term{x1}}
	ruleToMem(Predicate{"p",1}, Rule{[]Term{x1}, []Term{qx1, rx1}})
	x2 := &Var{"X"}
	ux2 := Compound_Term{Predicate{"u",1}, []Term{x2}}
	ruleToMem(Predicate{"p",1}, Rule{[]Term{x2}, []Term{ux2}})
	x3 := &Var{"X"}
	sx3 := Compound_Term{Predicate{"s",1}, []Term{x3}}
	ruleToMem(Predicate{"q",1}, Rule{[]Term{x3}, []Term{sx3}})
	ruleToMem(Predicate{"r",1}, Rule{[]Term{Atom{"a"}}, []Term{}})
	ruleToMem(Predicate{"r",1}, Rule{[]Term{Atom{"b"}}, []Term{}})
	ruleToMem(Predicate{"s",1}, Rule{[]Term{Atom{"a"}}, []Term{}})
	ruleToMem(Predicate{"s",1}, Rule{[]Term{Atom{"b"}}, []Term{}})
	ruleToMem(Predicate{"s",1}, Rule{[]Term{Atom{"c"}}, []Term{}})
	ruleToMem(Predicate{"u",1}, Rule{[]Term{Atom{"d"}}, []Term{}})
	x := &Var{"X"}
	px := Compound_Term{Predicate{"p",1}, []Term{x}}
	query := []Term{px}
	*/
	
	ruleToMem(Predicate{"int",1}, Rule{[]Term{Atom{"0"}}, []Term{}})
	m := &Var{"N"}
	sm := Compound_Term{Predicate{"s",1}, []Term{m}}
	im := Compound_Term{Predicate{"int",1}, []Term{m}}
	ruleToMem(Predicate{"int",1}, Rule{[]Term{sm}, []Term{im}})
	x := &Var{"X"}
	query := []Term{Compound_Term{Predicate{"int",1}, []Term{x}}}
	
	
	return query
}