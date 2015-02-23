
package prolog

import (
	t "GOPPL/types"
)

// newest item in stack at termlist[len-1] (->)
type Stack_Item struct {
	termlist t.Terms
	aliases  t.Alias
}

func InitStack(query t.Terms) Stack_Item {
	no_alias := make(t.Alias)
	return Stack_Item{query, no_alias}
}

func DFS(stack_item Stack_Item, answer chan t.Alias) {
	// for now, assume no parallellism
	// and only compound terms (no =, is etc)
	
	terms, aliases := stack_item.termlist, stack_item.aliases
	if len(terms) == 0 {
		answer <- aliases
		close(answer)
		return
	}
	tt, terms := terms[len(terms)-1], terms[:len(terms)-1]
	
	//Compound_Term assumption :
	term := tt.(t.Compound_Term)
	rules, contains := memory[term.Pred]
	if !contains {
		close(answer)
		return
	} else {
		for _,rule_template := range rules {
			rule := callRule(rule_template)
			new_terms := terms
			new_alias := make(t.Alias)
			scope := []*t.Var{}
			for k,v := range aliases {
				new_alias[k] = v
				scope = append(scope, k)
			}
			scope = append(arrangeVarsByDepth(scope), varsInTermArgs(term.GetArgs())...)
			unifies, al := t.Unify(term.GetArgs(), rule.Head, new_alias)
			if !unifies {
				continue
			}
			clash := t.UpdateAlias(new_alias, al)
			if clash { 
				continue 
			}			
			for i := len(rule.Body)-1; i >= 0; i-- {
				new_terms = append(new_terms, rule.Body[i])
			}
			si := Stack_Item{new_terms, new_alias}
			rec_answer := make(chan t.Alias)
			go DFS(si, rec_answer)
			for a := range rec_answer {
				a = cleanUpVarsOutOfScope(a, scope)
				answer <- a
			}
		}
		close(answer)
	}
}

//TODO: for efficiency, let rule templates use Var instead of *Var ?
func callRule(rule t.Rule) t.Rule {
	var_alias := make(t.Alias)
	head, body := t.Terms{}, t.Terms{}
	for _, term := range rule.Head {
		vt, var_alias := createVars(term, var_alias)
		var_alias = var_alias
		head = append(head, vt)
	}
	for _, term := range rule.Body {
		vt, var_alias := createVars(term, var_alias)
		var_alias = var_alias
		body = append(body, vt)
	}
	return t.Rule{head, body}
}

func createVars(term t.Term, va t.Alias) (t.Term, t.Alias) {
	switch term.(type) {
	case *t.Var:
		v := term.(*t.Var)
		value, renamed := va[v]
		if renamed {
			return value, va
		}
		newv := &t.Var{v.Name}	
		va[v] = newv
		return newv, va
	case t.Compound:
		renamed_args := t.Terms{}
		c := term.(t.Compound)
		for _, ot := range c.GetArgs() {
			vt, va := createVars(ot, va)
			va = va
			renamed_args = append(renamed_args, vt)
		}
		var newc t.Term
		switch c.(type) {
		case t.Compound_Term:
			newc = t.Compound_Term{c.GetPredicate(), renamed_args}
		case t.List:
			newc = t.List{t.Compound_Term{c.GetPredicate(), renamed_args}}
		}
		return newc, va
	}
	return term, va
}

func cleanUpVarsOutOfScope(to_clean t.Alias, scope []*t.Var) t.Alias {

	clean := make(t.Alias)
	for _, v := range scope {
		temp := v
		Loop: for {
			value, _ := to_clean[temp]
			switch value.(type) {
			case *t.Var:
				temp = value.(*t.Var)
			case t.Atom:
				clean[v] = value
				break Loop
			case t.Compound:
				compound := rec_substitute(value.(t.Compound), to_clean, scope)
				switch compound.(type) {
				case t.List:
					clean[v] = compound.(t.List)
				case t.Compound_Term:
					clean[v] = compound.(t.Compound_Term)
				}
				break Loop
			}
		}
	}
	return clean
}

func rec_substitute(c t.Compound, a t.Alias, scope []*t.Var) t.Compound {
	
	sub_args := t.Terms{}
	for _,term := range c.GetArgs() {
		switch term.(type){
		case t.Atom:
			sub_args = append(sub_args, term)
		case *t.Var:
			v := term.(*t.Var)
			v1, ok := a[v]
			if inScope(v, scope) || !ok {
				sub_args = append(sub_args, v)
			} else {	//var not in scope but bound in a
				switch v1.(type) {
				case t.Compound:
					sub_c := rec_substitute(v1.(t.Compound), a, scope)
					switch sub_c.(type) {
					case t.List:
						sub_args = append(sub_args, sub_c.(t.List))
					case t.Compound_Term:
						sub_args = append(sub_args, sub_c.(t.Compound_Term))
					}
				default:
					sub_args = append(sub_args, v1)
				}
			}
		case t.Compound:
			sub_c := rec_substitute(term.(t.Compound), a, scope)
			switch sub_c.(type) {
			case t.List:
				sub_args = append(sub_args, sub_c.(t.List))
			case t.Compound_Term:
				sub_args = append(sub_args, sub_c.(t.Compound_Term))
			}
		}
	}
	switch c.(type) {
	case t.List:
		return t.List{t.Compound_Term{c.GetPredicate(), sub_args}}
	}
	return t.Compound_Term{c.GetPredicate(), sub_args}
}

func varsInTermArgs(terms t.Terms) []*t.Var {
	vars := []*t.Var{}
	for _,term := range terms {
		switch term.(type) {
		case *t.Var:
			vars = append(vars, term.(*t.Var))
		case t.Compound_Term:
			vars = append(vars, varsInTermArgs(term.(t.Compound_Term).GetArgs())...)
		}
	}
	return vars
}

// TODO: arrange vars in scope by depth, then range from high depth to low
func arrangeVarsByDepth(scope []*t.Var) []*t.Var {
	return scope
}

func inScope(v *t.Var, scope []*t.Var) bool {
	for _, value := range scope {
		if value == v {
			return true
		}
	}
	return false
}