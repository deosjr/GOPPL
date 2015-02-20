
package prolog

import "fmt"

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
			for k,v := range aliases {
				new_alias[k] = v
			}
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
			scope := make(Alias)
			for k,v := range aliases {
				scope[k] = v
			}
			updateAlias(scope, varsInTermArgs(term.args))
			for i := len(rule.body)-1; i >= 0; i-- {
				new_terms = append(new_terms, rule.body[i])
			}
			si := Stack_Item{new_terms, new_alias}
			rec_answer := make(chan Alias)
			go DFS(si, rec_answer)
			for a := range rec_answer {
				fmt.Println("ANSWER", a, scope)
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

func cleanUpVarsOutOfScope(to_clean Alias, scope Alias) Alias {

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
				clean[k] = rec_substitute(value.(Compound_Term), to_clean, scope)
				break Loop
			}
		}
	}
	fmt.Println("CLEAN", clean)
	return clean
}

func varsInTermArgs(terms Terms) Alias {
	vars := make(Alias)
	for _,t := range terms {
		switch t.(type) {
		case *Var:
			vars[t.(*Var)] = Atom{"NIL"}
		case Compound_Term:
			updateAlias(vars, varsInTermArgs(t.(Compound_Term).args))
		}
	}
	return vars
}