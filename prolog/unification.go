
package prolog

func unify(args1 []Term, args2 []Term, aliases Bindings) (unified bool, newalias Bindings) {

	newalias = make(Bindings)
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

func (a Atom) UnifyWith(t Term, alias Bindings) (unified bool, newalias Bindings) {
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

func (v *Var) UnifyWith(t Term, a Bindings) (unified bool, newalias Bindings) {
	// already unified
	if bound, contains := a[v]; contains {
		return bound.UnifyWith(t, a)
	}
	// TODO: anonymous vars are all the same atm. Give them an identifier?
	if v.Name[0] == '_' {
		return true, a
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
	newalias = make(Bindings)
	newalias[v] = t
	return true, newalias
}

func (c Compound_Term) UnifyWith(t Term, alias Bindings) (unified bool, newalias Bindings) {
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

func (l List) UnifyWith(t Term, alias Bindings) (unified bool, newalias Bindings) {
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

func (v VarTemplate) UnifyWith(t Term, alias Bindings) (bool, Bindings) {
	return false, nil
}

func updateAlias(aliases Bindings, updates Bindings) (clash bool) {

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