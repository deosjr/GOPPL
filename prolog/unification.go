
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
		unifies, al := args1[i].UnifyWith(args2[i], aliases)
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

func (a Atom) UnifyWith(t Term, alias Alias) (unified bool, newalias Alias) {
	switch t.(type){
	case Atom:
		if a.Value == t.(Atom).Value {
			return true, newalias
		}
	case *Var:
		return t.UnifyWith(a, alias)
	}
	return false, nil
}

func (v *Var) UnifyWith(t Term, alias Alias) (unified bool, newalias Alias) {
	// already unified
	if bound, contains := alias[v]; contains {
		return bound.UnifyWith(t, alias)
	}
	switch t.(type){
	case Atom:
		// The RESERVED atom never unifies
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

func (c Compound_Term) UnifyWith(t Term, alias Alias) (unified bool, newalias Alias) {
	switch t.(type){
	case Compound_Term:
		ct := t.(Compound_Term)
		if c.GetPredicate() == ct.GetPredicate() {
			return unify(c.GetArgs(), ct.GetArgs(), alias)
		}
	case *Var:
		return t.UnifyWith(c, alias)
	}
	return false, nil
}

func (l List) UnifyWith(t Term, alias Alias) (unified bool, newalias Alias) {
	switch t.(type){
	case List:
		// This workaround is needed because RESERVED doesn't unify
		if l.compareTo(Empty_List) && t.(List).compareTo(Empty_List) {
			return true, alias
		}
		return unify(l.GetArgs(), t.(List).GetArgs(), alias)
	case *Var:
		return t.UnifyWith(l, alias)
	}
	return false, nil
}

func (v VarTemplate) UnifyWith(t Term, alias Alias) (bool, Alias) {
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