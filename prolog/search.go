
package prolog

//import "fmt"

type Data map[Predicate][]Rule

var Memory Data = make(Data)
var Extralogical = make( map[Predicate] func(Terms, Bindings) Bindings )

type searchnode struct {
	children []*searchnode
	term Compound_Term
	stack Terms 			// newest item in stack at end: stack[len-1] (->)
	answers []Bindings
	alias Bindings
	scope []*Var
	rules []Rule
	f func(Terms, Bindings) Bindings
	start bool
}

// returns nil if nothing remains to be explored
func newNode(terms Terms, alias Bindings) *searchnode {
	//fmt.Println("NODE: ", terms)
	if len(terms) == 0 {
		return nil
	}
	term := terms[len(terms)-1].(Compound_Term)
	scope := []*Var{}
	for k,_ := range alias {
		scope = append(scope, k)
	}
	scope = append(scope, VarsInTermArgs(term.GetArgs())...)
	return &searchnode{[]*searchnode{}, term, terms, []Bindings{}, alias, scope, []Rule{}, nil, true}
}

func EmptyDFS(query Terms) *searchnode {
	no_alias := make(Bindings)
	return BoundDFS(query, no_alias)
}

func BoundDFS(query Terms, alias Bindings) *searchnode {
	startnode := newNode(query, alias)
	return startnode
}

func (node *searchnode) GetAnswer() Bindings {
	if node.start {
		return node.startDFS()
	}
	return node.continueDFS()
}

func (node *searchnode) getTerms() Terms {
	return node.stack[:len(node.stack)-1]
}

func (node *searchnode) startDFS() Bindings {
	node.start = false
	rules, contains := Memory[node.term.Pred]
	if !contains {
		if f, ok := Extralogical[node.term.Pred]; ok {
			node.f = f
		} else {
			return nil
		}
	}
	node.rules = rules
	return node.continueDFS()
}

func (node *searchnode) continueDFS() Bindings {
	for len(node.children) != 0 {
		answer := node.getAnswerFromChild(node.children[0])
		if answer != nil {
			return answer
		}
	}
	if node.f != nil {
		return node.exploreFunction()
	}
	for len(node.rules) != 0 {
		answer := node.exploreRules()
		if answer != nil {
			return answer
		}
	}
	return nil
}

func (node *searchnode) exploreFunction() Bindings {
	new_alias := node.prepareExplore()
	alias := node.f(node.term.Args, node.alias)
	if alias == nil {
		return nil
	}
	return node.exploreFurther(new_alias, alias, node.getTerms())
}

func (node *searchnode) exploreRules() Bindings {
	if len(node.rules) == 0 {
		return nil
	}
	rule_template := node.rules[0]
	node.rules = node.rules[1:]
	rule := callRule(rule_template)
	//fmt.Println("CALL", node.term, rule)
	new_alias := node.prepareExplore()
	unifies, alias := unify(node.term.GetArgs(), rule.Head, new_alias)
	//fmt.Println("UNIFIES", unifies, term.GetArgs(), rule.Head, new_alias, al)
	if !unifies {
		return nil
	}
	return node.exploreFurther(new_alias, alias, appendNewTerms(node.getTerms(), rule.Body))
}

func (node *searchnode) prepareExplore() Bindings {
	new_alias := make(Bindings)
	for k,v := range node.alias {
		new_alias[k] = v
	}
	return new_alias
}

func (node *searchnode) exploreFurther(new_alias Bindings, alias Bindings, newterms Terms) Bindings {
	clash := UpdateAlias(new_alias, alias)
	if clash { 
		return nil
	}
	newnode := newNode(newterms, new_alias)
	if newnode == nil {
		return cleanUpVarsOutOfScope(new_alias, node.scope)
	}
	node.children = append(node.children, newnode)
	return node.getAnswerFromChild(newnode)
}

func (node *searchnode) getAnswerFromChild(child *searchnode) Bindings {
	answer := child.GetAnswer()
	if answer == nil {
		node.children = node.children[1:]
		return nil
	} else {
		return cleanUpVarsOutOfScope(answer, node.scope)
	}
}

func appendNewTerms(old Terms, new Terms) Terms {
	terms := make(Terms, len(old))
	copy(terms, old)
	for i := len(new)-1; i >= 0; i-- {
		terms = append(terms, new[i])
	}
	return terms
}

func callRule(rule Rule) Rule {
	var_alias := make(tempBindings)
	head, body := Terms{}, Terms{}
	for _, term := range rule.Head {
		vt, var_alias := term.CreateVars(var_alias)
		var_alias = var_alias
		head = append(head, vt)
	}
	for _, term := range rule.Body {
		vt, var_alias := term.CreateVars(var_alias)
		var_alias = var_alias
		body = append(body, vt)
	}
	return Rule{head, body}
}

func (a Atom) CreateVars(va tempBindings) (Term, tempBindings) {
	return a, va
}

func (i Int) CreateVars(va tempBindings) (Term, tempBindings) {
	return i, va
}

func (v *Var) CreateVars(va tempBindings) (Term, tempBindings) {
	return v, va
}

func (v VarTemplate) CreateVars(va tempBindings) (Term, tempBindings) {
	value, renamed := va[v]
	if renamed {
		return value, va
	}
	newv := &Var{v.Name}	
	va[v] = newv
	return newv, va
}

func (c Compound_Term) CreateVars(va tempBindings) (Term, tempBindings) {
	renamed_args := Terms{}
	for _, ot := range c.GetArgs() {
		vt, va := ot.CreateVars( va)
		va = va
		renamed_args = append(renamed_args, vt)
	}
	return Compound_Term{c.GetPredicate(), renamed_args}, va
}

func (n Nil) CreateVars(va tempBindings) (Term, tempBindings) {
	return n, va
}

func (c Cons) CreateVars(va tempBindings) (Term, tempBindings) {
	renamed_args := Terms{}
	for _, ot := range c.GetArgs() {
		vt, va := ot.CreateVars( va)
		va = va
		renamed_args = append(renamed_args, vt)
	}
	return Cons{Compound_Term{c.GetPredicate(), renamed_args}, renamed_args[0], renamed_args[1]}, va
}

func cleanUpVarsOutOfScope(to_clean Bindings, scope []*Var) Bindings {

	clean := make(Bindings)
	for _, v := range scope {	
		clean[v] = v.SubstituteVars(to_clean)
	}
	return clean
}

func VarsInTermArgs(terms Terms) []*Var {
	vars := []*Var{}
	for _, t := range terms {
		switch term := t.(type) {
		case *Var:
			vars = append(vars, term)
		case Cons:
			vars = append(vars, VarsInTermArgs(term.GetArgs())...)
		case Compound_Term:
			vars = append(vars, VarsInTermArgs(term.GetArgs())...)
		}
	}
	return vars
}