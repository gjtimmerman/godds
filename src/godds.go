package main

import (
	"fmt"
	"godds/interop"
	"slices"
)

func convertRankToChar(r int32) uint8 {
	switch r {
	case 2, 3, 4, 5, 6, 7, 8, 9:
		return uint8('0' + r)
	case 10:
		return 'T'
	case 11:
		return 'J'
	case 12:
		return 'Q'
	case 13:
		return 'K'
	case 14:
		return 'A'
	}
	return 0
}

func convertSuitToChar(s int32) uint8 {
	switch s {
	case 0:
		return 'S'
	case 1:
		return 'H'
	case 2:
		return 'D'
	case 3:
		return 'C'
	}
	return 0
}

func removeCard(remaining *[4][]byte, player int, suit int32, rank int32) {
	var suitCards int
	var start int
	for range suit {
		suitCards = slices.Index(remaining[player][start:], '.') + 1
		start += suitCards
	}
	cardPosition := slices.Index(remaining[player][start:], convertRankToChar(rank))
	cardPosition += start
	oldhand := remaining[player]
	remaining[player] = remaining[player][:cardPosition]
	remaining[player] = append(remaining[player], oldhand[cardPosition+1:]...)
}

func evaluateWinner(currentTrickSuit []int32, currentTrickRank []int32, winnerSoFar int, trumpSuit int32, suit int32, rank int32, currentTrickSequence int) int {
	if suit == trumpSuit {
		if currentTrickSuit[winnerSoFar] == trumpSuit {
			if rank > currentTrickRank[winnerSoFar] {
				return currentTrickSequence
			} else {
				return winnerSoFar
			}
		} else {
			return currentTrickSequence
		}
	} else {
		if currentTrickSuit[winnerSoFar] == trumpSuit {
			return winnerSoFar
		} else {
			if suit == currentTrickSuit[winnerSoFar] {
				if rank > currentTrickRank[winnerSoFar] {
					return currentTrickSequence
				} else {
					return winnerSoFar
				}
			} else {
				return winnerSoFar
			}
		}
	}
}

func main() {
	var remaincards []byte = []byte("N:T642.K9.AKT62.82 Q87.T7632.J83.74 AJ3.J85.Q7.KQJT5 K95.AQ4.954.A963")
	var remCards [4][]byte
	var space1 = 1
	var space2 = 0
	for i := range 4 {
		if space2 < len(remaincards) {
			space2 := slices.Index(remaincards[space1+1:], ' ')
			if space2 == -1 {
				space2 = len(remaincards)
			} else {
				space2 += space1 + 1
			}
			remCards[i] = remaincards[space1+1 : space2]
			//			fmt.Printf("%s\n", remCards[i])
			space1 = space2
		}
	}

	var deal interop.DealPBN
	deal.First = interop.NORTH
	deal.Trump = interop.SPADES
	ctx := interop.CreateSolverContext()

	for tr := range 13 {
		winnerSoFar := -1
		for t := range 4 {
			fmt.Printf("Turn: %d, %d\n", tr, t)
			copy(deal.RemainCards[:], "N:")
			nextCard := 2
			for i := range 4 {
				copy(deal.RemainCards[nextCard:], remCards[i])
				nextCard += len(remCards[i])
				deal.RemainCards[nextCard] = ' '
				nextCard++
			}

			for nextCard--; nextCard < 80; nextCard++ {
				deal.RemainCards[nextCard] = 0
			}
			//			fmt.Printf("%s\n", deal.RemainCards)
			var tricks interop.FutureTricks
			ret := interop.SolveBoardPbn(ctx, &deal, 1, 3, 0, &tricks)
			if ret != interop.RETURN_NO_FAULT {
				fmt.Printf("Error: %d", int32(ret))
				return
			}
			println(tricks.Score[0])
			//			for i := 0; i < int(tricks.Cards); i++ {
			fmt.Printf("%c", convertSuitToChar(tricks.Suit[0]))
			fmt.Printf("%c\n", convertRankToChar(tricks.Rank[0]))
			//			}
			if t != 3 {
				deal.CurrentTrickSuit[t] = tricks.Suit[0]
				deal.CurrentTrickRank[t] = tricks.Rank[0]
			}
			if winnerSoFar == -1 {
				winnerSoFar = 0
			} else {
				winnerSoFar = evaluateWinner(deal.CurrentTrickSuit[:], deal.CurrentTrickRank[:], winnerSoFar, deal.Trump, tricks.Suit[0], tricks.Rank[0], t)
			}
			removeCard(&remCards, (int(deal.First)+t)%4, tricks.Suit[0], tricks.Rank[0])
		}
		//		println(winnerSoFar)
		deal.First = (int32(winnerSoFar) + deal.First) % 4
		//		copy(deal.CurrentTrickSuit[:], []int32{0, 0, 0})
		//		copy(deal.CurrentTrickRank[:], []int32{0, 0, 0})
	}
}
