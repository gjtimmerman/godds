package main

import (
	"fmt"
	"godds/interop"
	"slices"
)

const (
	NS = iota
	EW
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

func convertHandToChar(h int) uint8 {
	switch h {
	case interop.NORTH:
		return 'N'
	case interop.EAST:
		return 'E'
	case interop.SOUTH:
		return 'S'
	case interop.WEST:
		return 'W'
	}
	return 0
}

type Card struct {
	Suit int32
	Rank int32
}

type Trick struct {
	Lead        int
	Winner      int
	PlayedCards [4]Card
}

type PlayedDeal struct {
	TricksWon [2]int
	Tricks    [13]Trick
}

type SolveDealError struct {
	Message string
}

func (sde SolveDealError) Error() string {
	return sde.Message
}

func removeCard(remaining *[interop.NUMBEROFPLAYERS][]byte, player int, suit int32, rank int32) {
	var suitCards int
	var start int
	for range suit {
		suitCards = slices.Index(remaining[player][start:], '.') + 1
		start += suitCards
	}
	cardPosition := slices.Index(remaining[player][start:], convertRankToChar(rank))
	cardPosition += start
	copy(remaining[player][cardPosition:], remaining[player][cardPosition+1:])
	remaining[player] = remaining[player][:len(remaining[player])-1]
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

func splitdeal(remaincards []byte) [4][]byte {
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
			space1 = space2
		}
	}
	return remCards
}

func solvedeal(ctx uintptr, first int32, trump int32, remaincards []byte) (playedDeal PlayedDeal, err error) {
	remCards := splitdeal(remaincards)
	var deal interop.DealPBN
	deal.Trump = trump
	deal.First = first
	for tr := range interop.NUMBEROFTRICKS {
		winnerSoFar := -1
		for t := range interop.NUMBEROFPLAYERS {
			copy(deal.RemainCards[:], "N:")
			nextCard := 2
			for i := range interop.NUMBEROFPLAYERS {
				copy(deal.RemainCards[nextCard:], remCards[i])
				nextCard += len(remCards[i])
				deal.RemainCards[nextCard] = ' '
				nextCard++
			}

			for nextCard--; nextCard < interop.MAX_REMAIN_CARDS; nextCard++ {
				deal.RemainCards[nextCard] = 0
			}
			var tricks interop.FutureTricks
			ret := interop.SolveBoardPbn(ctx, &deal, 1, 3, 0, &tricks)
			if ret != interop.RETURN_NO_FAULT {
				err = fmt.Errorf("Error: %d", int32(ret))
				return
			}
			playedDeal.Tricks[tr].Lead = int(deal.First)
			playedDeal.Tricks[tr].PlayedCards[t].Suit = tricks.Suit[0]
			playedDeal.Tricks[tr].PlayedCards[t].Rank = tricks.Rank[0]
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
		deal.First = (int32(winnerSoFar) + deal.First) % 4
		playedDeal.Tricks[tr].Winner = int(deal.First)
		if deal.First == interop.NORTH || deal.First == interop.SOUTH {
			playedDeal.TricksWon[NS]++
		} else {
			playedDeal.TricksWon[EW]++
		}
	}
	return
}

func printsolveddeal(playedDeal PlayedDeal) {
	fmt.Printf("Tricks won by NS: %d, EW: %d\n", playedDeal.TricksWon[NS], playedDeal.TricksWon[EW])
	for i, t := range playedDeal.Tricks {
		fmt.Printf("%d: Lead: %c, Winner: %c\n", i, convertHandToChar(t.Lead), convertHandToChar(t.Winner))
		for _, c := range t.PlayedCards {
			fmt.Printf("%c%c ", convertSuitToChar(c.Suit), convertRankToChar(c.Rank))
		}
		fmt.Println()
	}
}

func main() {
	ctx := interop.CreateSolverContext()
	var remaincards []byte = []byte("N:T642.K9.AKT62.82 Q87.T7632.J83.74 AJ3.J85.Q7.KQJT5 K95.AQ4.954.A963")
	playedDeal, err := solvedeal(ctx, interop.EAST, interop.SPADES, remaincards)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		printsolveddeal(playedDeal)
	}
}
