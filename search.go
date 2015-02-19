
package prolog

//import "fmt"

// newest item in stack at termlist[len-1] (->)
type Stack_Item struct {
	termlist []Term
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
		for _,rule := range rules {
			//fmt.Println("RULE", term.pred.functor, rule.head, rule.body)
			new_terms := terms
			new_alias := make(Alias)
			for k,v := range aliases {
				new_alias[k] = v
			}
			unifies, al := unify(term.args, rule.head, new_alias)
			//fmt.Println("UNIFIES?", unifies, term.args, rule.head)
			if !unifies {
				continue
			}
			//fmt.Println("UPDATE", new_alias, al)
			clash := update_alias(new_alias, al)
			//fmt.Println("CLASH?", clash)
			if clash { continue }
			for i := len(rule.body)-1; i >= 0; i-- {
				new_terms = append(new_terms, rule.body[i])
			}
			si := Stack_Item{new_terms, new_alias}
			rec_answer := make(chan Alias)
			go DFS(si, rec_answer)
			for a := range rec_answer {
				//fmt.Println("ANSWER",term,a)
				a = clean_up_vars_out_of_scope(a, new_alias)
				answer <- a
			}
		}
		close(answer)
	}
}