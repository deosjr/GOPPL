
package prolog

import "errors"

type Data map[Predicate][]Rule

var Memory Data = make(Data)

var Notification error = errors.New("notify")

type searchnode struct {
	Answer chan result
	notify chan bool
	children []*searchnode
	// newest item in stack at end: stack[len-1] (->)
	stack Terms
}

func (node *searchnode) Notify() {
	if len(node.children) > 0 {
		node.children[0].Notify()
	} else {
		node.notify <- true
	}
}

func (node *searchnode) closeAll(cont bool) {
	for _, c := range node.children {
		close(c.notify)
	}
	if cont && len(node.notify) > 0 {
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
	A Alias
	Err error
}

func StartDFS(query Terms) *searchnode {
	no_alias := make(Alias)
	startnode := newNode(query)
	go startnode.dfs(no_alias)
	startnode.Notify()
	return startnode
}

func (node *searchnode) dfs(aliases Alias) {
	
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
		node.closeAll(false)
		return
	}
	node.exploreRules(rules, term, terms, aliases)
	node.closeAll(true)
}

func (node *searchnode) exploreRules(rules []Rule, term Compound_Term, terms Terms, aliases Alias) {
	for _, rule_template := range rules {
		<- node.notify
		rule := callRule(rule_template)
		new_alias := make(Alias)
		scope := []*Var{}
		for k,v := range aliases {
			new_alias[k] = v
			scope = append(scope, k)
		}
		scope = append(scope, varsInTermArgs(term.GetArgs())...)
		unifies, al := unify(term.GetArgs(), rule.Head, new_alias)
		if !unifies {
			node.notify <- true
			continue
		}
		clash := updateAlias(new_alias, al)
		if clash { 
			node.notify <- true
			continue 
		}
		newnode := newNode(appendNewTerms(terms, rule.Body))
		node.children = append(node.children, newnode)
		go newnode.dfs(new_alias)
		node.awaitAnswers(newnode, scope)
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
	succes := false
	for r := range child.Answer {
		succes = true
		a, err := r.A, r.Err
		if err == Notification {
			node.notify <- true
			break
		}
		a = cleanUpVarsOutOfScope(a, scope)
		node.Answer <- result{a, err}
	}
	node.children = node.children[1:]
	close(child.notify)
	if !succes {
		node.notify <- true
	}
}

func callRule(rule Rule) Rule {
	var_alias := make(map[VarTemplate]Term)
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

func (a Atom) CreateVars(va map[VarTemplate]Term) (Term, map[VarTemplate]Term) {
	return a, va
}

func (v *Var) CreateVars(va map[VarTemplate]Term) (Term, map[VarTemplate]Term) {
	return v, va
}

func (v VarTemplate) CreateVars(va map[VarTemplate]Term) (Term, map[VarTemplate]Term) {
	value, renamed := va[v]
	if renamed {
		return value, va
	}
	newv := &Var{v.Name}	
	va[v] = newv
	return newv, va
}

func (c Compound_Term) CreateVars(va map[VarTemplate]Term) (Term, map[VarTemplate]Term) {
	renamed_args := Terms{}
	for _, ot := range c.GetArgs() {
		vt, va := ot.CreateVars( va)
		va = va
		renamed_args = append(renamed_args, vt)
	}
	return Compound_Term{c.GetPredicate(), renamed_args}, va
}

func (l List) CreateVars(va map[VarTemplate]Term) (Term, map[VarTemplate]Term) {
	renamed_args := Terms{}
	for _, ot := range l.GetArgs() {
		vt, va := ot.CreateVars( va)
		va = va
		renamed_args = append(renamed_args, vt)
	}
	return List{Compound_Term{l.GetPredicate(), renamed_args}}, va
}

func cleanUpVarsOutOfScope(to_clean Alias, scope []*Var) Alias {

	clean := make(Alias)
	for _, v := range scope {	
		temp := v
		Loop: for {
			value, ok := to_clean[temp]
			if !ok {
				break
			}
			switch value.(type) {
			case *Var:
				temp = value.(*Var)
			case Atom:
				clean[v] = value
				break Loop
			case Compound:
				clean[v] = value.(Compound).substituteVars(to_clean, scope)
				break Loop
			}
		}
	}
	return clean
}

func (a Atom) substituteVars(al Alias, scope []*Var) Term {
	return a
}

func (v VarTemplate) substituteVars(a Alias, scope []*Var) Term {
	return v
}

func (v *Var) substituteVars(a Alias, scope []*Var) Term {
	v1, ok := a[v]
	if inScope(v, scope) || !ok {
		return v
	}
	//var not in scope but bound in a
	switch v1.(type) {
	case Compound:
		return v1.(Compound).substituteVars(a, scope)
	}
	return v1
}

func (c Compound_Term) substituteVars(a Alias, scope []*Var) Term {
	
	sub_args := Terms{}
	for _,term := range c.GetArgs() {
		sub := term.substituteVars(a, scope)
		sub_args = append(sub_args, sub)
	}
	return Compound_Term{c.GetPredicate(), sub_args}
}

func (l List) substituteVars(a Alias, scope []*Var) Term {
	
	sub_args := Terms{}
	for _,term := range l.GetArgs() {
		sub := term.substituteVars(a, scope)
		sub_args = append(sub_args, sub)
	}
	return List{Compound_Term{l.GetPredicate(), sub_args}}
}

func varsInTermArgs(terms Terms) []*Var {
	vars := []*Var{}
	for _,term := range terms {
		switch term.(type) {
		case *Var:
			vars = append(vars, term.(*Var))
		case Compound_Term:
			vars = append(vars, varsInTermArgs(term.(Compound_Term).GetArgs())...)
		}
	}
	return vars
}

func inScope(v *Var, scope []*Var) bool {
	for _, value := range scope {
		if value == v {
			return true
		}
	}
	return false
}