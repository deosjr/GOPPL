
package prolog

import (
	"errors"
)

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
		clash := UpdateAlias(newalias, al)
		if clash {
			return false, nil
		}
	}
	return true, newalias
}

func (a Atom) UnifyWith(t Term, alias Bindings) (unified bool, newalias Bindings) {
	switch t.(type){
	case Atomic:
		if a.Value() == t.(Atomic).Value() {
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

func (c Cons) UnifyWith(t Term, alias Bindings) (unified bool, newalias Bindings) {
	switch t.(type){
	case Cons:
		return unify(c.GetArgs(), t.(List).GetArgs(), alias)
	case *Var:
		return t.UnifyWith(c, alias)
	}
	return false, nil
}

func (n Nil) UnifyWith(t Term, alias Bindings) (unified bool, newalias Bindings) {
	switch t.(type){
	case Nil:
		return true, newalias
	case *Var:
		return t.UnifyWith(n, alias)
	}
	return false, nil
}

func (v VarTemplate) UnifyWith(t Term, alias Bindings) (bool, Bindings) {
	return false, nil
}

func UpdateAlias(aliases Bindings, updates Bindings) (clash bool) {

	LOOP: for k, uv := range updates {
		if av, ok := aliases[k]; ok {
			switch av.(type) {
			case *Var:
				break
			default:
				// context is (alias/k) union (updates/k)
				// substituting both terms wrt the union of bindings
				// should be possible.
				context := make(Bindings)
				for kk, vv := range aliases {
					if k != kk {
						context[kk] = vv
					}
				}
				for kk, vv := range updates {
					if k != kk {
						context[kk] = vv
					}
				}
				subav := av.SubstituteVars(context)
				subuv := uv.SubstituteVars(context)
				if !subuv.compareTo(subav) {
					return true
				}
				aliases[k] = subav
				continue LOOP
			}
		}
		aliases[k] = uv
	}
	return false
}

func (a Atom) SubstituteVars(al Bindings) Term {
	return a
}

func (v VarTemplate) SubstituteVars(a Bindings) Term {
	return v
}

func (v *Var) SubstituteVars(a Bindings) Term {
	v1, ok := a[v]
	if !ok {
		return v
	}
	return v1.SubstituteVars(a)
}

func (c Compound_Term) SubstituteVars(a Bindings) Term {
	
	sub_args := Terms{}
	for _,term := range c.GetArgs() {
		sub := term.SubstituteVars(a)
		sub_args = append(sub_args, sub)
	}
	return Compound_Term{c.GetPredicate(), sub_args}
}

func (n Nil) SubstituteVars(a Bindings) Term {
	return n
}

func (c Cons) SubstituteVars(a Bindings) Term {
	
	sub_args := Terms{}
	for _,term := range c.GetArgs() {
		sub := term.SubstituteVars(a)
		sub_args = append(sub_args, sub)
	}
	return Cons{Compound_Term{c.GetPredicate(), sub_args}, sub_args[0], sub_args[1]}
}

var InstantiationError error = errors.New("arguments insufficiently instantiated")

func Evaluate(t Term, a Bindings) (int, error) {
	switch t.(type) {
	case Int:
		return t.(Int).IntValue(), nil
	case *Var:
		value, ok := a[t.(*Var)]
		if ok {
			return Evaluate(value, a)
		} 
		return 0, InstantiationError
	case Compound_Term:
		ct := t.(Compound_Term)
		if ct.Pred.Arity != 2 {
			return 0, nil
		}
		v1, err1 := Evaluate(ct.Args[0], a)
		v2, err2 := Evaluate(ct.Args[1], a)
		if err1 != nil {
			return 0, err1
		}
		if err2 != nil {
			return 0, err2
		}
		switch ct.Pred.Functor {
		case "+":
			return v1 + v2, nil
		case "-":
			return v1 - v2, nil
		case "*":
			return v1 * v2, nil
		case "/":
			return v1 / v2, nil
		}
	}
	return 0, nil
}