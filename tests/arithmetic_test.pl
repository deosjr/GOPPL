split(L,0,[],L).
split([X|Xs],N,[X|Ys],Zs) :- 	
	>(N,0), 
	is(N1, -(N,1)), 
	split(Xs,N1,Ys,Zs).