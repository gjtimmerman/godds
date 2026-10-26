package interop

import (
	"syscall"
	"unsafe"
)

const MAXNOOFBOARDS int32 = 200

const NUMBEROFTRICKS int32 = 13

const RETURN_NO_FAULT uintptr = 1

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
	RemainCards      [80]byte
}

type FutureTricks struct {
	Nodes  int32
	Cards  int32
	Suit   [NUMBEROFTRICKS]int32
	Rank   [NUMBEROFTRICKS]int32
	Equals [NUMBEROFTRICKS]int32
	Score  [NUMBEROFTRICKS]int32
}

var ddsdll, err1 = syscall.LoadLibrary("DDS.dll")
var createcontext, err2 = syscall.GetProcAddress(ddsdll, "dds_create_solvercontext_default")
var solveboard, err3 = syscall.GetProcAddress(ddsdll, "dds_solve_board_pbn")

func CreateSolverContext() uintptr {
	ret1, _, _ := syscall.SyscallN(createcontext)
	return ret1
}

func SolveBoardPbn(ctx uintptr, deal *DealPBN, target int32, solutions int32, mode int32, tricks *FutureTricks) uintptr {
	ret1, _, _ := syscall.SyscallN(solveboard, ctx, uintptr(unsafe.Pointer(deal)), uintptr(target), uintptr(solutions), uintptr(mode), uintptr(unsafe.Pointer(tricks)))
	return ret1
}
