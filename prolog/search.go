
package prolog

//import "fmt"

// newest item in stack at termlist[len-1] (->)
type Stack_Item struct {
	termlist Terms
	aliases  Alias
}

func InitStack(query []Term) Stack_Item {
	no_alias := make(Alias)
	return Stack_Item{query, no_alias}
}

func DFS(stack_item Stack_Item, answer chan Alias) {
	// for now, assume no parallellism
	// and only compound terms (no =, is etc)
	
	terms, aliases := stack_item.termlist, stack_item.aliases
	//fmt.Println("CALL", terms, aliases)
	if len(terms) == 0 {
		answer <- aliases
		close(answer)
		return
	}
	t, terms := terms[len(terms)-1], terms[:len(terms)-1]
	
	//Compound_Term assumption :
	term := t.(Compound_Term)
	rules, contains := memory[term.pred]
	if !contains {
		close(answer)
		return
	} else {
		for _,rule_template := range rules {
			rule := callRule(rule_template)
			//fmt.Println("RULE", term.pred.functor, rule.head, rule.body)
			new_terms := terms
			new_alias := make(Alias)
			scope := []*Var{}
			for k,v := range aliases {
				new_alias[k] = v
				scope = append(scope, k)
			}
			scope = append(arrangeVarsByDepth(scope), varsInTermArgs(term.args)...)
			unifies, al := unify(term.args, rule.head, new_alias)
			//fmt.Println("UNIFIES?", unifies, term.args, rule.head, new_alias, al)
			if !unifies {
				continue
			}
			//fmt.Println("UPDATE", new_alias, al)
			clash := updateAlias(new_alias, al)
			//fmt.Println("CLASH?", clash)
			if clash { 
				continue 
			}			
			for i := len(rule.body)-1; i >= 0; i-- {
				new_terms = append(new_terms, rule.body[i])
			}
			si := Stack_Item{new_terms, new_alias}
			rec_answer := make(chan Alias)
			go DFS(si, rec_answer)
			for a := range rec_answer {
				//fmt.Println("ANSWER", a, scope)
				a = cleanUpVarsOutOfScope(a, scope)
				answer <- a
			}
		}
		close(answer)
	}
}

//TODO: for efficiency, let rule templates use Var instead of *Var ?
func callRule(rule Rule) Rule {
	var_alias := make(Alias)
	head, body := []Term{}, []Term{}
	for _, t := range rule.head {
		vt, var_alias := createVars(t, var_alias)
		var_alias = var_alias
		head = append(head, vt)
	}
	for _, t := range rule.body {
		vt, var_alias := createVars(t, var_alias)
		var_alias = var_alias
		body = append(body, vt)
	}
	return Rule{head, body}
}

func createVars(t Term, va Alias) (Term, Alias) {
	switch t.(type) {
	case *Var:
		v := t.(*Var)
		value, renamed := va[v]
		if renamed {
			return value, va
		}
		newv := &Var{v.name}	
		va[v] = newv
		return newv, va
	case Compound_Term:
		renamed_args := []Term{}
		c := t.(Compound_Term)
		for _, ot := range c.args {
			vt, va := createVars(ot, va)
			va = va
			renamed_args = append(renamed_args, vt)
		}
		newc := Compound_Term{c.pred, renamed_args}
		return newc, va
	}
	return t, va
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

func cleanUpVarsOutOfScope(to_clean Alias, scope []*Var) Alias {

	clean := make(Alias)
	for _, v := range scope {
		temp := v
		Loop: for {
			value, _ := to_clean[temp]
			switch value.(type) {
			case *Var:
				temp = value.(*Var)
			case Atom:
				clean[v] = value
				break Loop
			case Compound_Term:
				clean[v] = rec_substitute(value.(Compound_Term), to_clean, scope)
				break Loop
			}
		}
	}
	return clean
}

func rec_substitute(c Compound_Term, a Alias, scope []*Var) Compound_Term {
	
	sub_args := []Term{}
	for _,t := range c.args {
		switch t.(type){
		case Atom:
			sub_args = append(sub_args, t)
		case *Var:
			v := t.(*Var)
			v1, ok := a[v]
			if inScope(v, scope) || !ok {
				sub_args = append(sub_args, v)
			} else {	//var not in scope but bound in a
				switch v1.(type) {
				case Compound_Term:
					sub_c := rec_substitute(v1.(Compound_Term), a, scope)
					sub_args = append(sub_args, sub_c)
				default:
					sub_args = append(sub_args, v1)
				}
			}
		case Compound_Term:
			sub_c := rec_substitute(t.(Compound_Term), a, scope)
			sub_args = append(sub_args, sub_c)
		}
	}
	return Compound_Term{c.pred, sub_args}
}

func varsInTermArgs(terms Terms) []*Var {
	vars := []*Var{}
	for _,t := range terms {
		switch t.(type) {
		case *Var:
			vars = append(vars, t.(*Var))
		case Compound_Term:
			vars = append(vars, varsInTermArgs(t.(Compound_Term).args)...)
		}
	}
	return vars
}

// TODO: arrange vars in scope by depth, then range from high depth to low
func arrangeVarsByDepth(scope []*Var) []*Var {
	return scope
}

func inScope(v *Var, scope []*Var) bool {
	for _, value := range scope {
		if value == v {
			return true
		}
	}
	return false
}