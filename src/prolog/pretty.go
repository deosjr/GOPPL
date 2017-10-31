
package prolog

func (a Atom) String() string{ return a.Value()}
func (i Int) String() string{ return i.Value()}

func (v *Var) String() string { return v.Name }
func (v VarTemplate) String() string { return v.Name }

func (c Compound_Term) String() string{ 
	if c.Pred.Arity == 0 {
		return c.Pred.Functor
	}
	s := c.Pred.Functor + "("
	for i,t := range c.Args {
		if i == c.Pred.Arity-1 {
			s += t.String()
			break
		}
		s += t.String() + ","
	}
	return s + ")"
}

func (c Cons) String() string { 
	tail := c.Tail()
	switch tail.(type){
	case List:
		ltail := tail.(List)
		if ltail.isEmpty() {
			return "[" + c.Head().String() + "]"
		}
		rec := ltail.String()
		return  "[" + c.Head().String() + "," + rec[1:len(rec)-1] + "]"
	}
	// TODO: doesnt take into account var X halfway string thats otherwise grounded
	return  "[" + c.Head().String() + "|" + tail.String() + "]"
}

func (n Nil) String() string {
	return "[]"
}

// For debugging only:

func (tlist Terms) String() string {
	s := "<"
	for _,t := range tlist {
		s = s + t.String() + " "
	}
	return s + ">"
} 

func (a Bindings) String() string {
	s := "{"
	for k,v := range a {
		s = s + k.String() + ":" + v.String() + " "
	}
	return s + "}"
}