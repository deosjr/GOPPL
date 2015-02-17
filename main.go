package main

import (
	"fmt"
)

type Rule struct {
	head []Term
	body []Term
}

type Predicate struct {
	functor string
	arity   int
}

type Term interface {
	term_to_string() string
}

type Atom struct {
	value string
}

type Var struct {
	name string
}

type Compound_Term struct {
	pred Predicate
	args []Term
}

func (a Atom) term_to_string() string { return a.value }

func (v *Var) term_to_string() string { return v.name }

func (v *Var) String() string { return v.name }

func (c Compound_Term) term_to_string() string {
	s := c.pred.functor + "("
	for i, t := range c.args {
		if i == c.pred.arity-1 {
			s += t.term_to_string()
			break
		}
		s += t.term_to_string() + ","
	}
	return s + ")"
}

func print(m map[Predicate][]Rule) {
	for k, v := range m {
		for _, rule := range v {
			fmt.Printf("%s(", k.functor)
			for i, h := range rule.head {
				if i == k.arity-1 {
					fmt.Printf("%s)", h.term_to_string())
					break
				}
				fmt.Printf("%s,", h.term_to_string())
			}
			if len(rule.body) == 0 {
				fmt.Println(".")
			} else {
				fmt.Println(" :-")
				for i, b := range rule.body {
					if i == len(rule.body)-1 {
						fmt.Printf("\t%s.", b.term_to_string())
						break
					}
					fmt.Printf("\t%s,\n", b.term_to_string())
				}
				fmt.Println()
			}
		}
		fmt.Println()
	}
}

func ruleToMem(pred Predicate, r Rule) {
	if value, ok := memory[pred]; ok {
		memory[pred] = append(value, r)
	} else {
		memory[pred] = []Rule{r}
	}
}

func update_alias(aliases map[*Var]Term, updates map[*Var]Term) (clash bool) {

	for k, v := range updates {
		if av, ok := aliases[k]; ok {
			switch av.(type) {
			case *Var:
				break
			default:
				if av != v {
					return true
				}
			}
		}
		aliases[k] = v
	}
	return false
}

func unify(args1 []Term, args2 []Term, aliases map[*Var]Term) (unified bool, newalias map[*Var]Term) {

	newalias = make(map[*Var]Term)
	for k, v := range aliases {
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

func unify_term(term1 Term, term2 Term, aliases map[*Var]Term) (unified bool, newalias map[*Var]Term) {

	newalias = make(map[*Var]Term)

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

func rename_alias(alias map[*Var]Term, t1 *Var, t2 Term) (bool, map[*Var]Term) {

	newalias := make(map[*Var]Term)

	var temp Term = t1
Loop:
	for {
		value, contains := alias[temp.(*Var)]
		if !contains {
			break
		}
		switch value.(type) {
		case *Var:
			temp = value
		default:
			if value == t2 {
				newalias[t1] = t2
				return true, newalias
			} else {
				break Loop
			}
		}

	}
	return false, newalias
}

var memory map[Predicate][]Rule = make(map[Predicate][]Rule)

// newest item in stack at termlist[len-1] (->)
type Stack_Item struct {
	termlist []Term
	aliases  map[*Var]Term
}

func dfs(stack_item Stack_Item, answer chan map[*Var]Term) {
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
		for _, rule := range rules {
			//fmt.Println("RULE", term.pred.functor, rule.head, rule.body)
			new_terms := terms
			new_alias := make(map[*Var]Term)
			for k, v := range aliases {
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
			if clash {
				continue
			}
			for i := len(rule.body) - 1; i >= 0; i-- {
				new_terms = append(new_terms, rule.body[i])
			}
			si := Stack_Item{new_terms, new_alias}
			rec_answer := make(chan map[*Var]Term)
			go dfs(si, rec_answer)
			for a := range rec_answer {
				//fmt.Println("ANSWER",term,a)
				//TODO: simplify a (clean up vars not in scope)
				answer <- a
			}
		}
		close(answer)
	}
}

func main() {

	ruleToMem(Predicate{"sym", 1}, Rule{[]Term{Atom{"a"}}, []Term{}})
	ruleToMem(Predicate{"sym", 1}, Rule{[]Term{Atom{"b"}}, []Term{}})
	h1 := &Var{"H1"}
	h2 := &Var{"H2"}
	sym1 := Compound_Term{Predicate{"sym", 1}, []Term{h1}}
	sym2 := Compound_Term{Predicate{"sym", 1}, []Term{h2}}
	ruleToMem(Predicate{"hardcoded2", 2}, Rule{[]Term{h1, h2}, []Term{sym1, sym2}})
	x, y := &Var{"X"}, &Var{"Y"}
	h := Compound_Term{Predicate{"hardcoded2", 2}, []Term{x, y}}
	query := []Term{h}

	/**
	ruleToMem(Predicate{"p",1}, Rule{[]Term{Atom{"a"}}, []Term{}})
	x1 := &Var{"X1"}
	qx1 := Compound_Term{Predicate{"q",1}, []Term{x1}}
	rx1 := Compound_Term{Predicate{"r",1}, []Term{x1}}
	ruleToMem(Predicate{"p",1}, Rule{[]Term{x1}, []Term{qx1, rx1}})
	x2 := &Var{"X2"}
	ux2 := Compound_Term{Predicate{"u",1}, []Term{x2}}
	ruleToMem(Predicate{"p",1}, Rule{[]Term{x2}, []Term{ux2}})
	x3 := &Var{"X3"}
	sx3 := Compound_Term{Predicate{"s",1}, []Term{x3}}
	ruleToMem(Predicate{"q",1}, Rule{[]Term{x3}, []Term{sx3}})
	ruleToMem(Predicate{"r",1}, Rule{[]Term{Atom{"a"}}, []Term{}})
	ruleToMem(Predicate{"r",1}, Rule{[]Term{Atom{"b"}}, []Term{}})
	ruleToMem(Predicate{"s",1}, Rule{[]Term{Atom{"a"}}, []Term{}})
	ruleToMem(Predicate{"s",1}, Rule{[]Term{Atom{"b"}}, []Term{}})
	ruleToMem(Predicate{"s",1}, Rule{[]Term{Atom{"c"}}, []Term{}})
	ruleToMem(Predicate{"u",1}, Rule{[]Term{Atom{"d"}}, []Term{}})
	x := &Var{"X"}
	px := Compound_Term{Predicate{"p",1}, []Term{x}}
	query := []Term{px}
	*/

	//print(memory)

	no_alias := make(map[*Var]Term)
	answer := make(chan map[*Var]Term, 1)
	stack_item := Stack_Item{query, no_alias}

	go dfs(stack_item, answer)

	fmt.Printf("?- %s.\n", query[0].term_to_string())
	wait := false
	for alias := range answer {
		for k, v := range alias {
			fmt.Printf("%s = %s. ", k, v.term_to_string())
		}
		if wait {
			for {
				var response string
				fmt.Scanln(&response)
				if response == ";" {
					break
				}
				if response == "a" {
					wait = false
					break
				}
			}
		}
		fmt.Println()
	}
	fmt.Println("False.")

}
