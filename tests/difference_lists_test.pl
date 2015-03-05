pal(X) :- pal(X, []).

pal(A, A).
pal([0|A], A).
pal([1|A], A).
pal([0|A], C) :-
        pal(A, [0|C]).
pal([1|A], C) :-
        pal(A, [1|C]).

% pal --> [] doesnt parse because zero args compounds arent parsed yet!
pal(X) --> [].
pal(X) --> [0].
pal(X) --> [1].
pal(X) --> [0], pal, [0].
pal(X) --> [1], pal, [1].