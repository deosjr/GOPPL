% see http://homepage.cs.uiowa.edu/~fleck/dcgTrans.htm

pal(0) --> [].
pal(0) --> [0].
pal(0) --> [1].
pal(s(X)) --> [0], pal(X), [0].
pal(s(X)) --> [1], pal(X), [1].

pal --> [].
pal --> [0].
pal --> [1].
pal --> [0], pal, [0].
pal --> [1], pal, [1].

%pal(A, A).
%pal([0|A], A).
%pal([1|A], A).
%pal([0|A], C) :-
%        pal(A, [0|C]).
%pal([1|A], C) :-
%        pal(A, [1|C]).