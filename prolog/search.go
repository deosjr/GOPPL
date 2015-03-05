
package prolog

import "errors"
//import "fmt"

type Data map[Predicate][]Rule

var Memory Data = make(Data)
var Extralogical = make( map[Predicate] func(Terms, Bindings) Bindings )

var Notification error = errors.New("notify")

type searchnode struct {
	Answer chan result
	notify chan bool
	children []*searchnode
	// newest item in stack at end: stack[len-1] (->)
	stack Terms
}

func (node *searchnode) wait() {
	<- node.notify
}

func (node *searchnode) Notify() {
	if len(node.children) > 0 {
		node.children[0].Notify()
	} else {
		node.notify <- true
	}
}

func (node *searchnode) notifySelf() {
	node.notify <- true
}

func (node *searchnode) closeAll(pass_on bool) {
	for _, c := range node.children {
		close(c.notify)
	}
	if pass_on && len(node.notify) > 0 {
		<- node.notify
		node.Answer <- result{nil, Notification}
	}
	close(node.Answer)
}

func newNode(terms Terms) *searchnode {
	answer := make(chan result, 1)
	notify := make(chan bool, 1)
	return &searchnode{answer, notify, []*searchnode{}, terms}
}

type result struct {
	Alias Bindings
	Err error
}

func StartDFS(query Terms) *searchnode {
	no_alias := make(Bindings)
	return ContinueDFS(query, no_alias)
}

func ContinueDFS(query Terms, alias Bindings) *searchnode {
	startnode := newNode(query)
	go startnode.dfs(alias)
	startnode.Notify()
	return startnode
}

func (node *searchnode) dfs(aliases Bindings) {
	
	if len(node.stack) == 0 {
		node.Answer <- result{aliases, nil}
		node.closeAll(false)
		return
	}
	terms, t := node.stack[:len(node.stack)-1], node.stack[len(node.stack)-1]
	
	//Compound_Term assumption (TODO: check at parse?):
	term := t.(Compound_Term)
	rules, contains := Memory[term.Pred]
	if !contains {
		if f, ok := Extralogical[term.Pred]; ok {
			node.exploreFunction(f, term, terms, aliases)
		} else {
			node.closeAll(false)
			return
		}
	} else {
		node.exploreRules(rules, term, terms, aliases)
	}
	node.closeAll(true)
}

func prepareExplore(term Compound_Term, aliases Bindings) (Bindings, []*Var) {
	new_alias := make(Bindings)
	scope := []*Var{}
	for k,v := range aliases {
		new_alias[k] = v
		scope = append(scope, k)
	}
	scope = append(scope, VarsInTermArgs(term.GetArgs())...)
	return new_alias, scope
}

func (node *searchnode) exploreFurther(new_alias Bindings, al Bindings, scope []*Var, newterms Terms) {
	clash := UpdateAlias(new_alias, al)
	if clash { 
		node.notifySelf()
		return
	}
	newnode := newNode(newterms)
	node.children = append(node.children, newnode)
	go newnode.dfs(new_alias)
	node.awaitAnswers(newnode, scope)
}

func (node *searchnode) exploreFunction(f func(Terms, Bindings) Bindings, term Compound_Term, terms Terms, aliases Bindings) {
	node.wait()
	new_alias, scope := prepareExplore(term, aliases)
	al := f(term.Args, aliases)
	if al == nil {
		return
	}
	node.exploreFurther(new_alias, al, scope, terms)
}

func (node *searchnode) exploreRules(rules []Rule, term Compound_Term, terms Terms, aliases Bindings) {
	for _, rule_template := range rules {
		node.wait()
		rule := callRule(rule_template)
		//fmt.Println("CALL", rule)
		new_alias, scope := prepareExplore(term, aliases)
		unifies, al := unify(term.GetArgs(), rule.Head, new_alias)
		//fmt.Println("UNIFIES", unifies, term.GetArgs(), rule.Head, new_alias, al)
		if !unifies {
			node.notifySelf()
			continue
		}
		node.exploreFurther(new_alias, al, scope, appendNewTerms(terms, rule.Body))
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

func (node *searchnode) awaitAnswers(child *searchnode, scope []*Var) {
	child.Notify()
	found_nothing := true
	for r := range child.Answer {
		found_nothing = false
		a, err := r.Alias, r.Err
		if err == Notification {
			node.notifySelf()
			break
		}
		a = cleanUpVarsOutOfScope(a, scope)
		node.Answer <- result{a, err}
	}
	node.children = node.children[1:]
	close(child.notify)
	if found_nothing {
		node.notifySelf()
	}
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
	for _, term := range terms {
		switch term.(type) {
		case *Var:
			vars = append(vars, term.(*Var))
		case Cons:
			vars = append(vars, VarsInTermArgs(term.(Cons).GetArgs())...)
		case Compound_Term:
			vars = append(vars, VarsInTermArgs(term.(Compound_Term).GetArgs())...)
		}
	}
	return vars
}