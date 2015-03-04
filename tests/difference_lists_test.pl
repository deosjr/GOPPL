pal(X) :- pal(X, []).

pal(A, A).
pal([0|A], A).
pal([1|A], A).
pal([0|A], C) :-
        pal(A, [0|C]).
pal([1|A], C) :-
        pal(A, [1|C]).