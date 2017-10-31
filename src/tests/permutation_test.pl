sym(a).
sym(b).

mys(b).
mys(c).

hardcoded2(H1, H2) :-
	sym(H1), sym(H2).

updateAliasProblem(X) :-
	sym(X), mys(X).