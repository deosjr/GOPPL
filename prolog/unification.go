
package prolog

//import "fmt"

func unify(args1 []Term, args2 []Term, aliases Alias) (unified bool, newalias Alias) {

	newalias = make(Alias)
	for k,v := range aliases {
		newalias[k] = v
	}
	
	if len(args1) != len(args2) {
		return false, nil
	}
	
	for i := 0; i < len(args1); i++ {
		unifies, al := unify_term(args1[i], args2[i], aliases)
		if !unifies {
			//fmt.Println("TERMS DONT UNIFY")
			return false, nil
		}
		clash := updateAlias(newalias, al)
		if clash {
			//fmt.Println("CLASH FROM UNIFY", newalias, al)
			return false, nil
		}
	}
	return true, newalias
}

func unify_term(term1 Term, term2 Term, aliases Alias) (unified bool, newalias Alias) {

	newalias = make(Alias)
	// RESERVED is an ununifyable atom
	if atom1, ok := term1.(Atom); ok && atom1.Value == "RESERVED" {
		return false, nil
	}
	if atom2, ok := term2.(Atom); ok && atom2.Value == "RESERVED" {
		return false, nil
	}

	// unification of var1:
	if var1, ok := term1.(*Var); ok {
		// already unified
		if bound, contains := aliases[var1]; contains {
			return unify_term(bound, term2, aliases)
		// var1 and var2
		} else if var2, ok2 := term2.(*Var); ok2 {
			newalias[var1] = var2
			return true, newalias	
		// var1 and nonvar2
		} else {
			newalias[var1] = term2
			return true, newalias
		}
	// unification of var2
	} else if var2, ok := term2.(*Var); ok {
		// already unified
		if bound, contains := aliases[var2]; contains {
			return unify_term(term1, bound, aliases)
		// var2 and nonvar1
		} else {
			newalias[var2] = term1
			return true, newalias
		}
	// unification of two atoms:
	} else if atom1, ok1 := term1.(Atom); ok1 {
		if atom2, ok2 := term2.(Atom); ok2 {
			if atom1.Value == atom2.Value {
				return true, newalias
			}
		}
	// can't unify compound term with atom
	} else if _, ok2 := term2.(Atom); ok2 { 
		return false, nil
	// unification of two compound terms
	} else if c1, c2 := term1.(Compound), term2.(Compound); c1.GetPredicate() == c2.GetPredicate() {
		return unify(c1.GetArgs(), c2.GetArgs(), aliases)
	}
	return false, nil
}

func updateAlias(aliases Alias, updates Alias) (clash bool) {

	for k,v := range updates {
		if av, ok := aliases[k]; ok {
			switch av.(type) {
			case *Var:
				break
			default:
				if !av.compare_to(v) {
					return true
				}
			}
		}
		aliases[k] = v
	}
	return false
}