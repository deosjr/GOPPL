
package prolog

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
	compareTo(Term) bool
	substituteVars(Bindings) Term
}

// TODO: separate int Atom
type Atom struct {
	Value string
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

// TODO: distinction between ground and unground compound terms
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

func (a Atom) compareTo(t Term) bool {
	switch t.(type) {
	case Atom:
		return a == t
	}
	return false
}

func (v *Var) compareTo(t Term) bool {
	switch t.(type) {
	case *Var:
		return v == t
	}
	return false
}

func (c Compound_Term) compareTo(t Term) bool {
	switch t.(type) {
	case Compound:
		tc := t.(Compound)
		if c.GetPredicate() == tc.GetPredicate() {
			cargs, tcargs := c.GetArgs(), tc.GetArgs()
			for i:=0; i < len(cargs); i++ {
				if !cargs[i].compareTo(tcargs[i]) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (v VarTemplate) compareTo(t Term) bool {
	switch t.(type) {
	case VarTemplate:
		return v.Name == t.(VarTemplate).Name
	}
	return false
}