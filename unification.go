
package prolog

//import "fmt"

func update_alias(aliases Alias, updates Alias) (clash bool) {

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
		clash := update_alias(newalias, al)
		if clash {
			//fmt.Println("CLASH FROM UNIFY", newalias, al)
			return false, nil
		}
	}
	return true, newalias
}

func unify_term(term1 Term, term2 Term, aliases Alias) (unified bool, newalias Alias) {

	newalias = make(Alias)

	// unification of two atoms:
	if atom1, ok1 := term1.(Atom); ok1 {
		if atom2, ok2 := term2.(Atom); ok2 {
			if atom1.value == atom2.value {
				return true, newalias
			}
		}
	// unification of var1:
	} else if var1, ok := term1.(*Var); ok {
		// already unified
		if _, contains := aliases[var1]; contains {
			if renamed, newalias := rename_alias(aliases, var1, term2); renamed {
				return true, newalias
			}
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
		if _, contains := aliases[var2]; contains {
			if renamed, newalias := rename_alias(aliases, var2, term1); renamed {
				return true, newalias
			}
		// var2 and nonvar1
		} else {
			newalias[var2] = term1
			return true, newalias
		}
	// unification of two compound terms
	} else if c1, c2 := term1.(Compound_Term), term2.(Compound_Term); c1.pred == c2.pred {
		return unify(c1.args, c2.args, aliases)
	}
	return false, nil
}

func rename_alias(alias Alias, t1 *Var, t2 Term) (bool, Alias){

	newalias := make(Alias)
	
	var temp Term = t1
	Loop: for {
		value, contains := alias[temp.(*Var)]
		if !contains { break }
		switch value.(type) {
		case *Var:
			temp = value
		default:
			if value.compare_to(t2) {
				newalias[t1] = t2
				return true, newalias
			} else { break Loop }
		}
	}
	return false, newalias
}

func clean_up_vars_out_of_scope(to_clean Alias, scope Alias) Alias {

	clean := make(Alias)
	for k,_ := range scope {
		var temp Term = k
		Loop: for {
			value, _ := to_clean[temp.(*Var)]
			switch value.(type) {
			case *Var:
				temp = value
			case Atom:
				clean[k] = value
				break Loop
			case Compound_Term:
				//TODO: X = s(N), N = 0 --> X = s(0)
				clean[k] = rec_substitute(value.(Compound_Term), to_clean, scope)
				break Loop
			}
		}
	}
	return clean
}

func rec_substitute(c Compound_Term, a Alias, scope Alias) Compound_Term {
	
	sub_args := []Term{}
	for _,t := range c.args {
		switch t.(type){
		case Atom:
			sub_args = append(sub_args, t)
		case *Var:
			v := t.(*Var)
			v1, ok := a[v]
			_, in_scope := scope[v]
			if in_scope || !ok {
				sub_args = append(sub_args, v)
			} else {	//var not in scope but bound in a
				sub_args = append(sub_args, v1)
			}
		case Compound_Term:
			sub_c := rec_substitute(t.(Compound_Term), a, scope)
			sub_args = append(sub_args, sub_c)
		}
	}
	return Compound_Term{c.pred, sub_args}
}