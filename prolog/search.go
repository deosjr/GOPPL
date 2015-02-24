
package prolog

type Data map[Predicate][]Rule

var Memory Data = make(Data)

// newest item in stack at termlist[len-1] (->)
type Stack_Item struct {
	termlist Terms
	aliases  Alias
}

func InitStack(query Terms) Stack_Item {
	no_alias := make(Alias)
	return Stack_Item{query, no_alias}
}

func DFS(stack_item Stack_Item, answer chan Alias) {
	// for now, assume no parallellism
	
	terms, aliases := stack_item.termlist, stack_item.aliases
	if len(terms) == 0 {
		answer <- aliases
		close(answer)
		return
	}
	t, terms := terms[len(terms)-1], terms[:len(terms)-1]
	
	//Compound_Term assumption :
	term := t.(Compound_Term)
	rules, contains := Memory[term.Pred]
	if !contains {
		close(answer)
		return
	} else {
		for _,rule_template := range rules {
			rule := callRule(rule_template)
			new_terms := terms
			new_alias := make(Alias)
			scope := []*Var{}
			for k,v := range aliases {
				new_alias[k] = v
				scope = append(scope, k)
			}
			scope = append(scope, varsInTermArgs(term.GetArgs())...)
			unifies, al := unify(term.GetArgs(), rule.Head, new_alias)
			if !unifies {
				continue
			}
			clash := updateAlias(new_alias, al)
			if clash { 
				continue 
			}			
			for i := len(rule.Body)-1; i >= 0; i-- {
				new_terms = append(new_terms, rule.Body[i])
			}
			si := Stack_Item{new_terms, new_alias}
			rec_answer := make(chan Alias)
			go DFS(si, rec_answer)
			for a := range rec_answer {
				a = cleanUpVarsOutOfScope(a, scope)
				answer <- a
			}
		}
		close(answer)
	}
}

func callRule(rule Rule) Rule {
	var_alias := make(map[VarTemplate]Term)
	head, body := Terms{}, Terms{}
	for _, term := range rule.Head {
		vt, var_alias := CreateVars(term, var_alias)
		var_alias = var_alias
		head = append(head, vt)
	}
	for _, term := range rule.Body {
		vt, var_alias := CreateVars(term, var_alias)
		var_alias = var_alias
		body = append(body, vt)
	}
	return Rule{head, body}
}

func CreateVars(term Term, va map[VarTemplate]Term) (Term, map[VarTemplate]Term) {
	switch term.(type) {
	case VarTemplate:
		v := term.(VarTemplate)
		value, renamed := va[v]
		if renamed {
			return value, va
		}
		newv := &Var{v.Name}	
		va[v] = newv
		return newv, va
	case Compound:
		renamed_args := Terms{}
		c := term.(Compound)
		for _, ot := range c.GetArgs() {
			vt, va := CreateVars(ot, va)
			va = va
			renamed_args = append(renamed_args, vt)
		}
		var newc Term
		switch c.(type) {
		case Compound_Term:
			newc = Compound_Term{c.GetPredicate(), renamed_args}
		case List:
			newc = List{Compound_Term{c.GetPredicate(), renamed_args}}
		}
		return newc, va
	}
	return term, va
}

//TODO: Loops when memory is read from file?
func cleanUpVarsOutOfScope(to_clean Alias, scope []*Var) Alias {

	clean := make(Alias)
	for _, v := range scope {	
		temp := v
		Loop: for {
			value, ok := to_clean[temp]
			if !ok {
				break
			}
			switch value.(type) {
			case *Var:
				temp = value.(*Var)
			case Atom:
				clean[v] = value
				break Loop
			case Compound:
				compound := rec_substitute(value.(Compound), to_clean, scope)
				switch compound.(type) {
				case List:
					clean[v] = compound.(List)
				case Compound_Term:
					clean[v] = compound.(Compound_Term)
				}
				break Loop
			}
		}
	}
	return clean
}

func rec_substitute(c Compound, a Alias, scope []*Var) Compound {
	
	sub_args := Terms{}
	for _,term := range c.GetArgs() {
		switch term.(type){
		case Atom:
			sub_args = append(sub_args, term)
		case *Var:
			v := term.(*Var)
			v1, ok := a[v]
			if inScope(v, scope) || !ok {
				sub_args = append(sub_args, v)
			} else {	//var not in scope but bound in a
				switch v1.(type) {
				case Compound:
					sub_c := rec_substitute(v1.(Compound), a, scope)
					switch sub_c.(type) {
					case List:
						sub_args = append(sub_args, sub_c.(List))
					case Compound_Term:
						sub_args = append(sub_args, sub_c.(Compound_Term))
					}
				default:
					sub_args = append(sub_args, v1)
				}
			}
		case Compound:
			sub_c := rec_substitute(term.(Compound), a, scope)
			switch sub_c.(type) {
			case List:
				sub_args = append(sub_args, sub_c.(List))
			case Compound_Term:
				sub_args = append(sub_args, sub_c.(Compound_Term))
			}
		}
	}
	switch c.(type) {
	case List:
		return List{Compound_Term{c.GetPredicate(), sub_args}}
	}
	return Compound_Term{c.GetPredicate(), sub_args}
}

func varsInTermArgs(terms Terms) []*Var {
	vars := []*Var{}
	for _,term := range terms {
		switch term.(type) {
		case *Var:
			vars = append(vars, term.(*Var))
		case Compound_Term:
			vars = append(vars, varsInTermArgs(term.(Compound_Term).GetArgs())...)
		}
	}
	return vars
}

func inScope(v *Var, scope []*Var) bool {
	for _, value := range scope {
		if value == v {
			return true
		}
	}
	return false
}