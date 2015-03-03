cat([], L, L).
cat([H|T], L, [H|R]) :-
	cat(T, L, R).

length([], 0).
length([H|T], N) :-
	length(T, N1),
	is(N, +(N1, 1)).