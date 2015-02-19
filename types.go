
package prolog

type Rule struct {
	head []Term
	body []Term
}

type Predicate struct {
	functor string
	arity int
}

type Term interface {
	Term_to_string() string
	compare_to(Term) bool
	ground(Alias) bool
}

type Atom struct {
	value string
}

type Var struct {
	name string
}

// TODO: distinction between ground and unground compound terms
type Compound_Term struct {
	pred Predicate
	args []Term
}

func (a Atom) Term_to_string() string{ return a.value}

func (v *Var) Term_to_string() string{ return v.name}

func (v *Var) String() string { return v.name }

func (c Compound_Term) Term_to_string() string{ 
	s := c.pred.functor + "("
	for i,t := range c.args {
		if i == c.pred.arity-1 {
			s += t.Term_to_string()
			break
		}
		s += t.Term_to_string() + ","
	}
	return s + ")"
}

func (a Atom) compare_to(t Term) bool {
	switch t.(type) {
	case Atom:
		return a == t
	}
	return false
}

func (v *Var) compare_to(t Term) bool {
	switch t.(type) {
	case *Var:
		return v == t
	}
	return false
}

func (c Compound_Term) compare_to(t Term) bool {
	switch t.(type) {
	case Compound_Term:
		tc := t.(Compound_Term)
		if c.pred == tc.pred {
			for i:=0; i < len(c.args); i++ {
				if !c.args[i].compare_to(tc.args[i]) {
					return false
				}
			}
			return true
		}
	}
	return false
}

func (a Atom) ground(alias Alias) bool {
	return true
}

// Grounded Var is bound to Atom or Compound_Term
// TODO: Groundedness of Compound doesnt matter right now. Should it?
func (v *Var) ground(alias Alias) bool {
	if value,contains := alias[v]; contains {
		_,ok := value.(*Var)
		return !ok
	}
	return false
}

func (c Compound_Term) ground(alias Alias) bool {
	for _,t := range c.args {
		if !t.ground(alias) {
			return false
		}
	}
	return true
}

type Alias map[*Var]Term