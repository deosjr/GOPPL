p(a).
p(X) :- q(X), r(X).
p(X) :- u(X).
	
q(X) :- s(X).

u(d).

r(a). r(b).
s(a). s(b). s(c).

q(1,_).
r(_,2).

p(X,Y) :- q(X,Y), r(X,Y).

test1(p(1,_)).
test2(p(_,2)).

test(EEN) :- 
	test1(p(Z1,X)), 
	test2(p(Y,Z2)), 
	=(X,2), =(Y,1), 
	=(p(Z1,X), p(Y,Z2)),
	=(EEN, 1).

