
package prolog

import "strconv"

type Rule struct {
	Head Terms
	Body Terms
}

type Predicate struct {
	Functor string
	Arity int
}

type Terms []Term

type Bindings map[*Var]Term
type tempBindings map[VarTemplate]Term

type Term interface {
	String() string
	UnifyWith(Term, Bindings) (bool, Bindings)
	CreateVars(tempBindings) (Term, tempBindings)
	equals(Term) bool
	SubstituteVars(Bindings) Term
}

type Atomic interface {
	Term
	Value() string
}

type Atom struct {
	value string
}

type Int struct {
	Atom
	intvalue int
}

func GetInt(i int) Int {
	return Int{Atom{strconv.Itoa(i)}, i}
}

func GetAtomic(s string) Atomic {
	if i, err := strconv.Atoi(s); err == nil {
		return Int{Atom{s}, i}
	}
	return Atom{s}
}

func (a Atom) Value() string {
	return a.value
}

func (i Int) IntValue() int {
	return i.intvalue
}

type Var struct {
	Name string
}

type VarTemplate struct {
	Name string
}

type Compound interface {
	Term
	GetPredicate() Predicate
	GetArgs() Terms
}

type Compound_Term struct {
	Pred Predicate
	Args Terms
}

func (c Compound_Term) GetPredicate() Predicate {
	return c.Pred
}
func (c Compound_Term) GetArgs() Terms {
	return c.Args
}

func (a Atom) equals(t Term) bool {
	switch t.(type) {
	case Atomic:
		return a == t
	}
	return false
}

func (v *Var) equals(t Term) bool {
	switch t.(type) {
	case *Var:
		return v == t
	}
	return false
}

func (c Compound_Term) equals(t Term) bool {
	switch t.(type) {
	case Compound:
		tc := t.(Compound)
		if c.GetPredicate() == tc.GetPredicate() {
			cargs, tcargs := c.GetArgs(), tc.GetArgs()
			for i:=0; i < len(cargs); i++ {
				if !cargs[i].equals(tcargs[i]) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (v VarTemplate) equals(t Term) bool {
	switch t.(type) {
	case VarTemplate:
		return v.Name == t.(VarTemplate).Name
	}
	return false
}