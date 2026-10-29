package interop

import (
	"syscall"
	"unsafe"
)

const NUMBEROFPLAYERS = 4

const MAXNOOFBOARDS int32 = 200

const NUMBEROFTRICKS int32 = 13

const RETURN_NO_FAULT uintptr = 1

const MAX_REMAIN_CARDS = 80

const NUMBEROFSTRAINS = 5

const (
	NORTH = iota
	EAST
	SOUTH
	WEST
)

const (
	SPADES = iota
	HEARTHS
	DIAMONDS
	CLUBS
	NOTRUMP
)

type DealPBN struct {
	Trump            int32
	First            int32
	CurrentTrickSuit [3]int32
	CurrentTrickRank [3]int32
	RemainCards      [MAX_REMAIN_CARDS]byte
}

type FutureTricks struct {
	Nodes  int32
	Cards  int32
	Suit   [NUMBEROFTRICKS]int32
	Rank   [NUMBEROFTRICKS]int32
	Equals [NUMBEROFTRICKS]int32
	Score  [NUMBEROFTRICKS]int32
}

type DdTableDealPBN struct {
	Cards [80]byte
}

type DdTableDealsPBN struct {
	No_of_tables int32
	Deals        [MAXNOOFBOARDS]DdTableDealPBN
}

type ParResults struct {
	/* index = 0 is NS view and index = 1
	   is EW view. By 'view' is here meant
	   which side that starts the bidding. */
	Par_score            [2][16]byte
	Par_contracts_string [2][128]byte
}

type AllParResults struct {
	Par_results [MAXNOOFBOARDS]ParResults
}

type DdTableResults struct {
	Res_table [NUMBEROFSTRAINS][NUMBEROFPLAYERS]int32
}

type DdTablesRes struct {
	No_of_boards int32
	Results      [MAXNOOFBOARDS * NUMBEROFSTRAINS]DdTableResults
}

var ddsdll, err1 = syscall.LoadLibrary("DDS.dll")
var createcontext, err2 = syscall.GetProcAddress(ddsdll, "dds_create_solvercontext_default")
var solveboard, err3 = syscall.GetProcAddress(ddsdll, "dds_solve_board_pbn")
var callalltables, err4 = syscall.GetProcAddress(ddsdll, "CalcAllTablesPBN")

func CreateSolverContext() uintptr {
	ret1, _, _ := syscall.SyscallN(createcontext)
	return ret1
}

func SolveBoardPbn(ctx uintptr, deal *DealPBN, target int32, solutions int32, mode int32, tricks *FutureTricks) uintptr {
	ret1, _, _ := syscall.SyscallN(solveboard, ctx, uintptr(unsafe.Pointer(deal)), uintptr(target), uintptr(solutions), uintptr(mode), uintptr(unsafe.Pointer(tricks)))
	return ret1
}

func CalcAllTablesPbn(deals *DdTableDealsPBN, mode int32, trumpfilter *int32, tableResults *DdTablesRes, parResults *AllParResults) uintptr {
	ret1, _, _ := syscall.SyscallN(callalltables, uintptr(unsafe.Pointer(deals)), uintptr(mode), uintptr(unsafe.Pointer(trumpfilter)), uintptr(unsafe.Pointer(tableResults)), uintptr(unsafe.Pointer(parResults)))
	return ret1
}
