package main

import "godds/interop"

func main() {
	var remaincards []byte = []byte("N:T642.K9.AKT62.82 Q87.T7632.J83.74 AJ3.J85.Q7.KQJT5 K95.AQ4.954.A963")
	var deal interop.DealPBN
	deal.First = interop.NORTH
	deal.Trump = interop.SPADES
	copy(deal.RemainCards[:], remaincards)
	var tricks interop.FutureTricks
	ctx := interop.CreateSolverContext()
	ret := interop.SolveBoardPbn(ctx, &deal, -1, 3, 0, &tricks)
	if ret == interop.RETURN_NO_FAULT {
		println("Success!")
	} else {
		println("Error")
		return
	}
	println(tricks.Cards)
}
