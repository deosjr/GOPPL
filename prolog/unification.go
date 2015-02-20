
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
			if atom1.value == atom2.value {
				return true, newalias
			}
		}
	// can't unify compound term with atom
	} else if _, ok2 := term2.(Atom); ok2 { 
		return false, nil
	// unification of two compound terms
	} else if c1, c2 := term1.(Compound_Term), term2.(Compound_Term); c1.pred == c2.pred {
		return unify(c1.args, c2.args, aliases)
	}
	return false, nil
}