
package prolog

func unify(args1 []Term, args2 []Term, aliases Alias) (unified bool, newalias Alias) {

	newalias = make(Alias)
	for k,v := range aliases {
		newalias[k] = v
	}
	
	if len(args1) != len(args2) {
		return false, nil
	}
	
	for i := 0; i < len(args1); i++ {
		unifies, al := args1[i].unifyWith(args2[i], aliases)
		if !unifies {
			return false, nil
		}
		clash := updateAlias(newalias, al)
		if clash {
			return false, nil
		}
	}
	return true, newalias
}

func (a Atom) unifyWith(t Term, alias Alias) (unified bool, newalias Alias) {
	switch t.(type){
	case Atom:
		if a.Value == t.(Atom).Value {
			return true, newalias
		}
	case *Var:
		return t.unifyWith(a, alias)
	}
	return false, nil
}

func (v *Var) unifyWith(t Term, alias Alias) (unified bool, newalias Alias) {
	// already unified
	if bound, contains := alias[v]; contains {
		return bound.unifyWith(t, alias)
	}
	switch t.(type){
	case Atom:
		if t.(Atom).Value == "RESERVED" {
			return false, nil
		}	
	case *Var:
		// TODO: nothing? you sure?
	}
	newalias = make(Alias)
	newalias[v] = t
	return true, newalias
}

func (c Compound_Term) unifyWith(t Term, alias Alias) (unified bool, newalias Alias) {
	switch t.(type){
	case Compound_Term:
		ct := t.(Compound_Term)
		if c.GetPredicate() == ct.GetPredicate() {
			return unify(c.GetArgs(), ct.GetArgs(), alias)
		}
	case *Var:
		return t.unifyWith(c, alias)
	}
	return false, nil
}

func (l List) unifyWith(t Term, alias Alias) (unified bool, newalias Alias) {
	switch t.(type){
	case List:
		return unify(l.GetArgs(), t.(List).GetArgs(), alias)
	case *Var:
		return t.unifyWith(l, alias)
	}
	return false, nil
}

func (v VarTemplate) unifyWith(t Term, alias Alias) (bool, Alias) {
	return false, nil
}

func updateAlias(aliases Alias, updates Alias) (clash bool) {

	for k,v := range updates {
		if av, ok := aliases[k]; ok {
			switch av.(type) {
			case *Var:
				break
			default:
				if !av.compareTo(v) {
					return true
				}
			}
		}
		aliases[k] = v
	}
	return false
}