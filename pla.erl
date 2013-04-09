-module(pla).

-compile([export_all]).

gen_side() ->
    A = rand_point(),
    B = rand_point(),
    fun(C) -> 
            gen_point(A, B, C)
    end.

gen_point({Ax, Ay}, {Bx, By}, {Cx, Cy} = C) ->
    T = (Bx - Ax)*(Cy - Ay) - (By - Ay)*(Cx - Ax),
    if 
        T >= 0 -> {C, 1};
        T < 0 -> {C, -1}
    end.

rand_point() -> {rnd(), rnd()}.
rnd() -> crypto:rand_uniform(-100000, 100001) / 100000.

sign(A) when A >= 0 -> 1;
sign(_) -> -1.

input(N) ->
    Gen = gen_side(),
    {Gen, [Gen(rand_point()) || _ <- lists:seq(1, N)]}.

run_over(Times, N) ->
    Counters = [begin
                    {Gen, Input} = input(N),
                    {NCorr, Weights} = correcting(Input),
                    PMiss = misclassify(Gen, Weights),
                    {NCorr, PMiss}
                end || _ <- lists:seq(1, Times)],
    {NumsCorr, ProbsMiss} = lists:unzip(Counters),
    {lists:sum(NumsCorr) / Times, lists:sum(ProbsMiss) / Times}.
    

correcting(Input) ->
    pla(Input, {0, 0, 0}).

misclassify(Gen, Weights) ->
    Point = rand_point(),
    {Point, Sign} = Gen(Point),
    case pla_test({Point, Sign}, Weights) of
        true -> 0;
        false -> 1
    end.

pla(Input, Weights) ->
    pla(Input, Weights, 0).

pla(Input, Weights, Num) ->
    Mistaken = lists:filter(fun(In) -> not pla_test(In, Weights) end, Input),
    case length(Mistaken) of
        0 -> {Num, Weights};
        N -> InR = lists:nth(crypto:rand_uniform(1, N+1), Mistaken),
             NewWeights = pla_corr(InR, Weights),
             pla(Input, NewWeights, Num+1)
    end.
    
pla_test({{Px, Py}, Test}, {Wo, Wx, Wy}) ->
    case sign(Wo + Px*Wx + Py*Wy) =:= Test of
        true -> true;
        false -> false
    end.

pla_corr({{Px, Py}, Test}, {Wo, Wx, Wy}) ->
    {Wo + 1*Test, Wx + Px*Test, Wy + Py*Test}.
