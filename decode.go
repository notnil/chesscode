package chesscode

import (
	"errors"
	"fmt"
	"reflect"
	"sort"

	"github.com/notnil/chess"
)

var (
	maxPieces = startingPositionPieceCount()
)

func Decode(board *chess.Board) (s string, err error) {
	m := board.SquareMap()
	sqs := allSquares()
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()
	for _, p := range pieceOrder {
		psqs := squaresForPiece(m, p)
		idxs := indexesForSquares(sqs, psqs)
		if len(idxs) > maxPieces[p] {
			return "", errors.New("chesscode: invalid number of " + p.String())
		}
		if len(idxs) == 0 {
			return s, nil
		}
		switch p {
		case chess.WhiteKing, chess.BlackKing, chess.WhiteQueen, chess.BlackQueen:
			code := idxs[0]
			s += charset[code : code+1]
		case chess.WhiteRook, chess.BlackRook, chess.WhiteBishop, chess.BlackBishop:
			combos := combo2(len(sqs))
			comboIdx := searchCombos(idxs, combos)
			sub := char2RevLookup[comboIdx]
			s += sub
		case chess.WhiteKnight:
			bSqs := squaresForPiece(m, chess.BlackKnight)
			bIdxs := indexesForSquares(sqs, bSqs)
			r := []combo2by2Result{}
			for _, idx := range idxs {
				r = append(r, combo2by2Result{idx: idx, piece: chess.WhiteKnight})
			}
			for _, idx := range bIdxs {
				r = append(r, combo2by2Result{idx: idx, piece: chess.BlackKnight})
			}
			sort.Slice(r, func(i, j int) bool {
				return r[i].idx < r[j].idx
			})
			comboIdx := search2by2Combos(r, combo2by2Lookup)
			sub := char3RevLookup[comboIdx]
			s += sub
		case chess.WhitePawn, chess.BlackPawn:
			leftSqs := leftSquares(sqs)
			leftIdxs := indexesForSquares(leftSqs, psqs)
			combos := combo4(len(leftSqs))
			comboIdx := searchCombos(leftIdxs, combos)
			sub := char2RevLookup[comboIdx]
			s += sub

			rightSqs := rightSqares(sqs)
			rightIdxs := indexesForSquares(rightSqs, psqs)
			combos = combo4(len(rightSqs))
			comboIdx = searchCombos(rightIdxs, combos)
			sub = char2RevLookup[comboIdx]
			s += sub
		}
		sort.Sort(sort.Reverse(sort.IntSlice(idxs)))
		for _, idx := range idxs {
			sqs = popSquare(sqs, idx)
		}
	}
	return s, nil
}

func searchCombos(idxs []int, combos [][]int) int {
	for i, c := range combos {
		if reflect.DeepEqual(idxs, c) {
			return i
		}
	}
	panic("unreachable")
}

func search2by2Combos(idxs []combo2by2Result, combos [][]combo2by2Result) int {
	for i, c := range combos {
		if reflect.DeepEqual(idxs, c) {
			return i
		}
	}
	panic("unreachable")
}

func squaresForPiece(m map[chess.Square]chess.Piece, p chess.Piece) []chess.Square {
	sqs := []chess.Square{}
	for sq, sqp := range m {
		if p == sqp {
			sqs = append(sqs, sq)
		}
	}
	return sqs
}

func indexesForSquares(sqs []chess.Square, p []chess.Square) []int {
	results := []int{}
	for _, sq := range p {
		idx := indexOfSquare(sqs, sq)
		if idx != -1 {
			results = append(results, idx)
		}
	}
	sort.IntSlice(results).Sort()
	return results
}

func indexOfSquare(sqs []chess.Square, sq chess.Square) int {
	for i, sq2 := range sqs {
		if sq2 == sq {
			return i
		}
	}
	return -1
}

func startingPositionPieceCount() map[chess.Piece]int {
	m := map[chess.Piece]int{}
	sqMap := chess.StartingPosition().Board().SquareMap()
	for _, p := range sqMap {
		m[p] += 1
	}
	return m
}
