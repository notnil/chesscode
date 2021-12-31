package chesscode

import (
	"fmt"
	"sort"
	"strings"

	"github.com/notnil/chess"
)

const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ. "

var (
	pieceOrder = []chess.Piece{
		chess.WhiteKing,
		chess.BlackKing,
		chess.WhiteQueen,
		chess.BlackQueen,
		chess.WhiteRook,
		chess.BlackRook,
		chess.WhiteBishop,
		chess.BlackBishop,
		chess.WhiteKnight,
		chess.BlackKnight,
		chess.WhitePawn,
		chess.BlackPawn,
	}
)

func Encode(s string) (*chess.Board, error) {
	s = formatString(s)
	if err := validInputString(s); err != nil {
		return nil, err
	}
	sqs := allSquares()
	sqMap := map[chess.Square]chess.Piece{}
	blackKnightUsed := []int{}
	for _, p := range pieceOrder {
		if len(s) == 0 {
			break
		}
		used := []int{}
		chars := 0
		switch p {
		case chess.WhiteKing, chess.BlackKing, chess.WhiteQueen, chess.BlackQueen:
			code := string(s[0])
			idx := strings.Index(charset, code)
			used = append(used, idx)
			chars = 1
		case chess.WhiteRook, chess.BlackRook, chess.WhiteBishop, chess.BlackBishop:
			codes := s[:2]
			idx := char2Lookup[codes]
			combos := combo2(len(sqs))
			used = combos[idx]
			chars = 2
		case chess.WhiteKnight:
			codes := s[:3]
			idx := char3Lookup[codes]
			combo := combo2by2Lookup[idx]
			for _, c := range combo {
				switch c.piece {
				case chess.WhiteKnight:
					used = append(used, c.idx)
				case chess.BlackKnight:
					blackKnightUsed = append(blackKnightUsed, c.idx-len(used))
				}
			}
			chars = 0
		case chess.BlackKnight:
			used = append(used, blackKnightUsed...)
			chars = 3
		case chess.WhitePawn, chess.BlackPawn:
			leftSqs := leftSquares(sqs)
			codes := s[:2]
			lookupIdx := char2Lookup[codes]
			combos := combo4(len(leftSqs))
			combo := combos[lookupIdx]
			for i, c := range combo {
				sq := leftSqs[c]
				combo[i] = indexOfSquare(sqs, sq)
			}
			used = append(used, combo...)

			rightSqs := rightSqares(sqs)
			rightCodes := s[2:4]
			lookupIdx = char2Lookup[rightCodes]
			combos = combo4(len(rightSqs))
			combo = combos[lookupIdx]
			for i, c := range combo {
				sq := rightSqs[c]
				combo[i] = indexOfSquare(sqs, sq)
			}
			used = append(used, combo...)
			chars = 4
		}
		sort.Sort(sort.Reverse(sort.IntSlice(used)))
		for _, idx := range used {
			sq := sqs[idx]
			sqs = popSquare(sqs, idx)
			sqMap[sq] = p
		}
		s = s[chars:]
	}
	return chess.NewBoard(sqMap), nil
}

func formatString(s string) string {
	s = strings.ToUpper(s)
	l := len(s)
	boundaries := []int{0, 1, 2, 3, 4, 6, 8, 10, 12, 15, 19, 23}
	for _, b := range boundaries {
		if l == b {
			return s
		}
		if l < b {
			s += strings.Repeat(" ", b-l)
			break
		}
	}
	return s
}

func validInputString(s string) error {
	if len(s) > 23 {
		return fmt.Errorf("chesscode: maximum encoding length is 23 characters but got %d %s", len(s), s)
	}
	for _, c := range s {
		idx := strings.Index(charset, string(c))
		if idx == -1 {
			return fmt.Errorf("chesscode: invalid character %s in input %s", string(c), s)
		}
	}
	return nil
}

func allSquares() []chess.Square {
	sqs := []chess.Square{}
	for i := 0; i < 64; i++ {
		sqs = append(sqs, chess.Square(i))
	}
	return sqs
}

func leftSquares(sqs []chess.Square) []chess.Square {
	return filterSqs(sqs, func(s chess.Square) bool {
		switch s.File() {
		case chess.FileA, chess.FileB, chess.FileC,
			chess.FileD:
			return true
		}
		return false
	})
}

func rightSqares(sqs []chess.Square) []chess.Square {
	return filterSqs(sqs, func(s chess.Square) bool {
		switch s.File() {
		case chess.FileE, chess.FileF, chess.FileG,
			chess.FileH:
			return true
		}
		return false
	})
}

func filterSqs(sqs []chess.Square, fn func(chess.Square) bool) []chess.Square {
	cp := []chess.Square{}
	for _, sq := range sqs {
		if fn(sq) {
			cp = append(cp, sq)
		}
	}
	return cp
}

func popSquare(sqs []chess.Square, idx int) []chess.Square {
	return append(sqs[:idx], sqs[idx+1:]...)
}
