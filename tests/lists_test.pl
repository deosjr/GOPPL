cat([], L, L).
cat([H|T], L, [H|R]) :-
	cat(T, L, R).