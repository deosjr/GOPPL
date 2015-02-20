
package prolog

//TODO: 'builtin' boolean field?
type Rule struct {
	head Terms
	body Terms
}

type Predicate struct {
	functor string
	arity int
}

type Terms []Term

type Alias map[*Var]Term

type Term interface {
	String() string
	compare_to(Term) bool
	ground(Alias) bool
}

//TODO: Anonymous variables as Atoms with _value?
type Atom struct {
	value string
}

//TODO: (optimalization) Rules in memory should use Vars,
// only make *Vars in search.createVars()
type Var struct {
	name string
}

// TODO: distinction between ground and unground compound terms
type Compound_Term struct {
	pred Predicate
	args Terms
}

//TODO: check Equaler interface!

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
